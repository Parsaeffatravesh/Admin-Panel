package middleware

import (
        "bytes"
        "context"
        "fmt"
        "net/http"
        "strings"
        "sync"
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
        var mu sync.RWMutex
        clients := make(map[string]*client)

        return func(next http.Handler) http.Handler {
                return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                        ip := r.RemoteAddr

                        mu.Lock()
                        c, exists := clients[ip]
                        if !exists {
                                clients[ip] = &client{count: 1, lastSeen: time.Now()}
                                mu.Unlock()
                                next.ServeHTTP(w, r)
                                return
                        }

                        if time.Since(c.lastSeen) > window {
                                c.count = 1
                                c.lastSeen = time.Now()
                                mu.Unlock()
                                next.ServeHTTP(w, r)
                                return
                        }

                        if c.count >= requests {
                                mu.Unlock()
                                utils.ErrorResponse(w, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests", nil)
                                return
                        }

                        c.count++
                        mu.Unlock()
                        next.ServeHTTP(w, r)
                })
        }
}

type cacheEntry struct {
        data      []byte
        expiresAt time.Time
        status    int
        headers   http.Header
}

type responseRecorder struct {
        http.ResponseWriter
        body       *bytes.Buffer
        statusCode int
}

func (r *responseRecorder) Write(b []byte) (int, error) {
        r.body.Write(b)
        return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) WriteHeader(code int) {
        r.statusCode = code
        r.ResponseWriter.WriteHeader(code)
}

func Cache(ttl time.Duration) func(http.Handler) http.Handler {
        var mu sync.RWMutex
        cache := make(map[string]*cacheEntry)

        return func(next http.Handler) http.Handler {
                return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                        if r.Method != http.MethodGet {
                                next.ServeHTTP(w, r)
                                return
                        }

                        cacheKey := r.URL.Path + "?" + r.URL.RawQuery

                        mu.RLock()
                        entry, exists := cache[cacheKey]
                        if exists && time.Now().Before(entry.expiresAt) {
                                for k, v := range entry.headers {
                                        w.Header()[k] = v
                                }
                                w.Header().Set("X-Cache", "HIT")
                                w.WriteHeader(entry.status)
                                w.Write(entry.data)
                                mu.RUnlock()
                                return
                        }
                        mu.RUnlock()

                        rec := &responseRecorder{
                                ResponseWriter: w,
                                body:           &bytes.Buffer{},
                                statusCode:     http.StatusOK,
                        }

                        next.ServeHTTP(rec, r)

                        if rec.statusCode >= 200 && rec.statusCode < 300 {
                                mu.Lock()
                                cache[cacheKey] = &cacheEntry{
                                        data:      rec.body.Bytes(),
                                        expiresAt: time.Now().Add(ttl),
                                        status:    rec.statusCode,
                                        headers:   w.Header().Clone(),
                                }
                                mu.Unlock()
                        }
                })
        }
}

func SecurityHeaders(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("X-Content-Type-Options", "nosniff")
                w.Header().Set("X-Frame-Options", "DENY")
                w.Header().Set("X-XSS-Protection", "1; mode=block")
                w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
                next.ServeHTTP(w, r)
        })
}

func CacheableResponse(maxAge int) func(http.Handler) http.Handler {
        return func(next http.Handler) http.Handler {
                return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                        if r.Method == http.MethodGet {
                                w.Header().Set("Cache-Control", fmt.Sprintf("private, max-age=%d, stale-while-revalidate=60", maxAge))
                                w.Header().Set("Vary", "Authorization, Accept-Encoding")
                        }
                        next.ServeHTTP(w, r)
                })
        }
}

func NoCacheResponse(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
                w.Header().Set("Pragma", "no-cache")
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
