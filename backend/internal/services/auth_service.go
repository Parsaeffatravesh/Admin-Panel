package services

import (
        "context"
        "errors"
        "sync"
        "time"

        "admin-panel/internal/config"
        "admin-panel/internal/models"
        "admin-panel/internal/repository"
        "admin-panel/internal/utils"

        "github.com/golang-jwt/jwt/v5"
        "github.com/google/uuid"
)

var (
        ErrInvalidCredentials = errors.New("invalid credentials")
        ErrUserNotFound       = errors.New("user not found")
        ErrUserInactive       = errors.New("user account is inactive")
        ErrTokenExpired       = errors.New("token has expired")
        ErrInvalidToken       = errors.New("invalid token")
)

type TokenClaims struct {
        UserID   uuid.UUID `json:"user_id"`
        TenantID uuid.UUID `json:"tenant_id"`
        Email    string    `json:"email"`
        jwt.RegisteredClaims
}

type AuthTokens struct {
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
        ExpiresIn    int64  `json:"expires_in"`
}

type AuthService struct {
        userRepo       *repository.UserRepository
        roleRepo       *repository.RoleRepository
        sessionRepo    *repository.SessionRepository
        auditRepo      *repository.AuditLogRepository
        jwtConfig      config.JWTConfig
        permissionCache sync.Map
}

func NewAuthService(
        userRepo *repository.UserRepository,
        roleRepo *repository.RoleRepository,
        sessionRepo *repository.SessionRepository,
        auditRepo *repository.AuditLogRepository,
        jwtConfig config.JWTConfig,
) *AuthService {
        return &AuthService{
                userRepo:    userRepo,
                roleRepo:    roleRepo,
                sessionRepo: sessionRepo,
                auditRepo:   auditRepo,
                jwtConfig:   jwtConfig,
        }
}

type LoginRequest struct {
        Email    string `json:"email" validate:"required,email"`
        Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
        User   *models.User `json:"user"`
        Tokens *AuthTokens  `json:"tokens"`
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*LoginResponse, error) {
        // Lowercase email before lookup
        email := strings.ToLower(req.Email)
        user, err := s.userRepo.GetByEmail(ctx, email)
        if err != nil {
                return nil, ErrInvalidCredentials
        }

        if user.Status != "active" {
                return nil, ErrUserInactive
        }

        if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
                return nil, ErrInvalidCredentials
        }

        tokens, err := s.generateTokens(user)
        if err != nil {
                return nil, err
        }

        session := &models.Session{
                ID:           uuid.New(),
                UserID:       user.ID,
                RefreshToken: tokens.RefreshToken,
                IPAddress:    ipAddress,
                UserAgent:    userAgent,
                ExpiresAt:    time.Now().Add(s.jwtConfig.RefreshTokenTTL),
                CreatedAt:    time.Now(),
        }

        if err := s.sessionRepo.Create(ctx, session); err != nil {
                return nil, err
        }

        if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
                return nil, err
        }

        s.auditRepo.Log(ctx, &models.AuditLog{
                ID:        uuid.New(),
                TenantID:  user.TenantID,
                UserID:    &user.ID,
                Action:    "login",
                Resource:  "auth",
                IPAddress: ipAddress,
                UserAgent: userAgent,
                CreatedAt: time.Now(),
        })

        return &LoginResponse{
                User:   user,
                Tokens: tokens,
        }, nil
}

func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID, ipAddress, userAgent string) error {
        user, err := s.userRepo.GetByID(ctx, userID)
        if err != nil {
                return err
        }

        if err := s.sessionRepo.DeleteByUserID(ctx, userID); err != nil {
                return err
        }

        s.auditRepo.Log(ctx, &models.AuditLog{
                ID:        uuid.New(),
                TenantID:  user.TenantID,
                UserID:    &userID,
                Action:    "logout",
                Resource:  "auth",
                IPAddress: ipAddress,
                UserAgent: userAgent,
                CreatedAt: time.Now(),
        })

        return nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*AuthTokens, error) {
        claims, err := s.validateToken(refreshToken)
        if err != nil {
                return nil, err
        }

        session, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
        if err != nil {
                return nil, ErrInvalidToken
        }

        if session.ExpiresAt.Before(time.Now()) {
                return nil, ErrTokenExpired
        }

        user, err := s.userRepo.GetByID(ctx, claims.UserID)
        if err != nil {
                return nil, ErrUserNotFound
        }

        tokens, err := s.generateTokens(user)
        if err != nil {
                return nil, err
        }

        session.RefreshToken = tokens.RefreshToken
        session.ExpiresAt = time.Now().Add(s.jwtConfig.RefreshTokenTTL)
        if err := s.sessionRepo.Update(ctx, session); err != nil {
                return nil, err
        }

        return tokens, nil
}

func (s *AuthService) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
        return s.validateToken(tokenString)
}

func (s *AuthService) HasPermission(ctx context.Context, userID uuid.UUID, resource, action string) bool {
        cacheKey := userID.String()
        if cached, ok := s.permissionCache.Load(cacheKey); ok {
                permissions := cached.([]string)
                permissionStr := resource + ":" + action
                for _, p := range permissions {
                        if p == permissionStr || p == resource+":*" || p == "*:*" {
                                return true
                        }
                }
                return false
        }

        permissions, err := s.roleRepo.GetUserPermissions(ctx, userID)
        if err != nil {
                return false
        }

        permStrs := make([]string, len(permissions))
        for i, p := range permissions {
                permStrs[i] = p.Resource + ":" + p.Action
        }
        s.permissionCache.Store(cacheKey, permStrs)

        permissionStr := resource + ":" + action
        for _, p := range permStrs {
                if p == permissionStr || p == resource+":*" || p == "*:*" {
                        return true
                }
        }
        return false
}

func (s *AuthService) InvalidatePermissionCache(userID uuid.UUID) {
        s.permissionCache.Delete(userID.String())
}

func (s *AuthService) generateTokens(user *models.User) (*AuthTokens, error) {
        now := time.Now()

        accessClaims := &TokenClaims{
                UserID:   user.ID,
                TenantID: user.TenantID,
                Email:    user.Email,
                RegisteredClaims: jwt.RegisteredClaims{
                        ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtConfig.AccessTokenTTL)),
                        IssuedAt:  jwt.NewNumericDate(now),
                        NotBefore: jwt.NewNumericDate(now),
                        Issuer:    "admin-panel",
                        Subject:   user.ID.String(),
                },
        }

        accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
        accessTokenString, err := accessToken.SignedString([]byte(s.jwtConfig.Secret))
        if err != nil {
                return nil, err
        }

        refreshClaims := &TokenClaims{
                UserID:   user.ID,
                TenantID: user.TenantID,
                Email:    user.Email,
                RegisteredClaims: jwt.RegisteredClaims{
                        ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtConfig.RefreshTokenTTL)),
                        IssuedAt:  jwt.NewNumericDate(now),
                        NotBefore: jwt.NewNumericDate(now),
                        Issuer:    "admin-panel",
                        Subject:   user.ID.String(),
                },
        }

        refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
        refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtConfig.Secret))
        if err != nil {
                return nil, err
        }

        return &AuthTokens{
                AccessToken:  accessTokenString,
                RefreshToken: refreshTokenString,
                ExpiresIn:    int64(s.jwtConfig.AccessTokenTTL.Seconds()),
        }, nil
}

func (s *AuthService) validateToken(tokenString string) (*TokenClaims, error) {
        token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, ErrInvalidToken
                }
                return []byte(s.jwtConfig.Secret), nil
        })

        if err != nil {
                if errors.Is(err, jwt.ErrTokenExpired) {
                        return nil, ErrTokenExpired
                }
                return nil, ErrInvalidToken
        }

        claims, ok := token.Claims.(*TokenClaims)
        if !ok || !token.Valid {
                return nil, ErrInvalidToken
        }

        return claims, nil
}
