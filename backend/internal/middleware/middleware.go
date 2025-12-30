package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"admin-panel/internal/services"
	"admin-panel/internal/utils"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type contextKey string

const (
	UserContextKey     contextKey = "user"
	TenantContextKey   contextKey = "tenant"
	RequestIDKey       contextKey = "request_id"
)

type AuthMiddleware struct {
	authService *services.AuthService
	logger      zerolog.Logger
}

func NewAuthMiddleware(authService *services.AuthService, logger zerolog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.Unauthorized(w, "Missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.Unauthorized(w, "Invalid authorization header format")
			return
		}

		claims, err := m.authService.ValidateAccessToken(parts[1])
		if err != nil {
			utils.Unauthorized(w, "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequirePermission(resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*services.TokenClaims)
			if !ok {
				utils.Unauthorized(w, "User not authenticated")
				return
			}

			hasPermission := m.authService.HasPermission(r.Context(), claims.UserID, resource, action)
			if !hasPermission {
				utils.Forbidden(w, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequestLogger(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

			defer func() {
				logger.Info().
					Str("request_id", requestID).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Int("status", ww.Status()).
					Int("bytes", ww.BytesWritten()).
					Dur("duration_ms", time.Since(start)).
					Str("remote_addr", r.RemoteAddr).
					Msg("request completed")
			}()

			w.Header().Set("X-Request-ID", requestID)
			next.ServeHTTP(ww, r.WithContext(ctx))
		})
	}
}

func RateLimiter(requests int, window time.Duration) func(http.Handler) http.Handler {
	type client struct {
		count    int
		lastSeen time.Time
	}
	clients := make(map[string]*client)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			c, exists := clients[ip]
			if !exists {
				clients[ip] = &client{count: 1, lastSeen: time.Now()}
				next.ServeHTTP(w, r)
				return
			}

			if time.Since(c.lastSeen) > window {
				c.count = 1
				c.lastSeen = time.Now()
				next.ServeHTTP(w, r)
				return
			}

			if c.count >= requests {
				utils.ErrorResponse(w, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests", nil)
				return
			}

			c.count++
			next.ServeHTTP(w, r)
		})
	}
}

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func GetUserFromContext(ctx context.Context) *services.TokenClaims {
	claims, ok := ctx.Value(UserContextKey).(*services.TokenClaims)
	if !ok {
		return nil
	}
	return claims
}

func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return ""
	}
	return requestID
}
