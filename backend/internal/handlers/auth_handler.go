package handlers

import (
        "encoding/json"
        "net/http"
        "os"

        "admin-panel/internal/middleware"
        "admin-panel/internal/services"
        "admin-panel/internal/utils"

        "github.com/go-playground/validator/v10"
)

func isProduction() bool {
        env := os.Getenv("APP_ENV")
        return env == "production"
}

type AuthHandler struct {
        authService *services.AuthService
        validate    *validator.Validate
}

func NewAuthHandler(authService *services.AuthService, validate *validator.Validate) *AuthHandler {
        return &AuthHandler{
                authService: authService,
                validate:    validate,
        }
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
        var req services.LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                utils.BadRequest(w, "Invalid request body", nil)
                return
        }

        if err := h.validate.Struct(&req); err != nil {
                details := make(map[string]string)
                for _, e := range err.(validator.ValidationErrors) {
                        details[e.Field()] = e.Tag()
                }
                utils.BadRequest(w, "Validation failed", details)
                return
        }

        ipAddress := r.RemoteAddr
        userAgent := r.UserAgent()

        resp, err := h.authService.Login(r.Context(), &req, ipAddress, userAgent)
        if err != nil {
                switch err {
                case services.ErrInvalidCredentials:
                        utils.Unauthorized(w, "Invalid email or password")
                case services.ErrUserInactive:
                        utils.Forbidden(w, "Account is inactive")
                default:
                        utils.InternalError(w, "Login failed")
                }
                return
        }

        secure := isProduction()
        
        http.SetCookie(w, &http.Cookie{
                Name:     "access_token",
                Value:    resp.Tokens.AccessToken,
                Path:     "/",
                MaxAge:   900,
                HttpOnly: true,
                Secure:   secure,
                SameSite: http.SameSiteLaxMode,
        })

        http.SetCookie(w, &http.Cookie{
                Name:     "refresh_token",
                Value:    resp.Tokens.RefreshToken,
                Path:     "/",
                MaxAge:   604800,
                HttpOnly: true,
                Secure:   secure,
                SameSite: http.SameSiteLaxMode,
        })

        utils.JSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
        claims := middleware.GetUserFromContext(r.Context())
        if claims == nil {
                utils.Unauthorized(w, "Not authenticated")
                return
        }

        ipAddress := r.RemoteAddr
        userAgent := r.UserAgent()

        if err := h.authService.Logout(r.Context(), claims.UserID, ipAddress, userAgent); err != nil {
                utils.InternalError(w, "Logout failed")
                return
        }

        http.SetCookie(w, &http.Cookie{
                Name:     "access_token",
                Value:    "",
                Path:     "/",
                MaxAge:   -1,
                HttpOnly: true,
        })

        http.SetCookie(w, &http.Cookie{
                Name:     "refresh_token",
                Value:    "",
                Path:     "/",
                MaxAge:   -1,
                HttpOnly: true,
        })

        utils.JSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

type RefreshRequest struct {
        RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
        var req RefreshRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                utils.BadRequest(w, "Invalid request body", nil)
                return
        }

        if err := h.validate.Struct(&req); err != nil {
                utils.BadRequest(w, "Refresh token is required", nil)
                return
        }

        tokens, err := h.authService.RefreshTokens(r.Context(), req.RefreshToken)
        if err != nil {
                switch err {
                case services.ErrTokenExpired:
                        utils.Unauthorized(w, "Refresh token has expired")
                case services.ErrInvalidToken:
                        utils.Unauthorized(w, "Invalid refresh token")
                default:
                        utils.InternalError(w, "Token refresh failed")
                }
                return
        }

        utils.JSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
        claims := middleware.GetUserFromContext(r.Context())
        if claims == nil {
                utils.Unauthorized(w, "Not authenticated")
                return
        }

        utils.JSON(w, http.StatusOK, map[string]interface{}{
                "user_id":   claims.UserID,
                "tenant_id": claims.TenantID,
                "email":     claims.Email,
        })
}
