package main

import (
        "context"
        "net/http"
        "os"
        "os/signal"
        "syscall"
        "time"

        "admin-panel/internal/config"
        "admin-panel/internal/database"
        "admin-panel/internal/handlers"
        "admin-panel/internal/middleware"
        "admin-panel/internal/repository"
        "admin-panel/internal/services"

        "github.com/go-chi/chi/v5"
        chimiddleware "github.com/go-chi/chi/v5/middleware"
        "github.com/go-chi/cors"
        "github.com/go-playground/validator/v10"
        "github.com/rs/zerolog"
)

func main() {
        cfg := config.Load()

        logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
        if cfg.App.Environment == "development" {
                logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
        }

        logger.Info().Str("environment", cfg.App.Environment).Msg("Starting admin panel server")

        db, err := database.NewPostgresPool(&cfg.Database)
        if err != nil {
                logger.Fatal().Err(err).Msg("Failed to connect to database")
        }
        defer db.Close()

        logger.Info().Msg("Connected to database")

        userRepo := repository.NewUserRepository(db)
        roleRepo := repository.NewRoleRepository(db)
        sessionRepo := repository.NewSessionRepository(db)
        auditRepo := repository.NewAuditLogRepository(db)

        authService := services.NewAuthService(userRepo, roleRepo, sessionRepo, auditRepo, cfg.JWT)
        userService := services.NewUserService(userRepo, roleRepo, auditRepo)
        roleService := services.NewRoleService(roleRepo, auditRepo)
        auditService := services.NewAuditService(auditRepo)
        dashboardService := services.NewDashboardService(userRepo, roleRepo, auditRepo)

        validate := validator.New()

        authHandler := handlers.NewAuthHandler(authService, validate)
        userHandler := handlers.NewUserHandler(userService, validate)
        roleHandler := handlers.NewRoleHandler(roleService, validate)
        auditHandler := handlers.NewAuditHandler(auditService)
        dashboardHandler := handlers.NewDashboardHandler(dashboardService)

        authMiddleware := middleware.NewAuthMiddleware(authService, logger)

        r := chi.NewRouter()

        r.Use(chimiddleware.RequestID)
        r.Use(chimiddleware.RealIP)
        r.Use(chimiddleware.Compress(5))
        r.Use(middleware.RequestLogger(logger))
        r.Use(chimiddleware.Recoverer)
        r.Use(middleware.SecurityHeaders)
        r.Use(cors.Handler(cors.Options{
                AllowedOrigins:   cfg.Server.AllowedOrigins,
                AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
                AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
                ExposedHeaders:   []string{"X-Request-ID"},
                AllowCredentials: len(cfg.Server.AllowedOrigins) > 0 && cfg.Server.AllowedOrigins[0] != "*",
                MaxAge:           300,
        }))

        r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
                w.Write([]byte(`{"status":"healthy"}`))
        })

        r.Route("/api/v1", func(r chi.Router) {
                r.Route("/auth", func(r chi.Router) {
                        r.Use(middleware.RateLimiter(10, time.Minute))
                        r.Post("/login", authHandler.Login)
                        r.Post("/refresh", authHandler.RefreshToken)

                        r.Group(func(r chi.Router) {
                                r.Use(authMiddleware.Authenticate)
                                r.Post("/logout", authHandler.Logout)
                                r.Get("/me", authHandler.Me)
                        })
                })

                r.Group(func(r chi.Router) {
                        r.Use(authMiddleware.Authenticate)

                        r.Route("/dashboard", func(r chi.Router) {
                                r.With(authMiddleware.RequirePermission("dashboard", "read")).Get("/stats", dashboardHandler.GetStats)
                        })

                        r.Route("/users", func(r chi.Router) {
                                r.With(authMiddleware.RequirePermission("users", "read")).Get("/", userHandler.List)
                                r.With(authMiddleware.RequirePermission("users", "create")).Post("/", userHandler.Create)
                                r.With(authMiddleware.RequirePermission("users", "read")).Get("/{id}", userHandler.Get)
                                r.With(authMiddleware.RequirePermission("users", "update")).Put("/{id}", userHandler.Update)
                                r.With(authMiddleware.RequirePermission("users", "delete")).Delete("/{id}", userHandler.Delete)
                                r.With(authMiddleware.RequirePermission("users", "update")).Post("/{id}/reset-password", userHandler.ResetPassword)
                                r.With(authMiddleware.RequirePermission("users", "read")).Get("/{id}/roles", userHandler.GetRoles)
                        })

                        r.Route("/roles", func(r chi.Router) {
                                r.With(authMiddleware.RequirePermission("roles", "read")).Get("/", roleHandler.List)
                                r.With(authMiddleware.RequirePermission("roles", "create")).Post("/", roleHandler.Create)
                                r.With(authMiddleware.RequirePermission("roles", "read")).Get("/{id}", roleHandler.Get)
                                r.With(authMiddleware.RequirePermission("roles", "update")).Put("/{id}", roleHandler.Update)
                                r.With(authMiddleware.RequirePermission("roles", "delete")).Delete("/{id}", roleHandler.Delete)
                                r.With(authMiddleware.RequirePermission("roles", "read")).Get("/{id}/permissions", roleHandler.GetPermissions)
                        })

                        r.Route("/permissions", func(r chi.Router) {
                                r.With(authMiddleware.RequirePermission("roles", "read")).Get("/", roleHandler.GetAllPermissions)
                        })

                        r.Route("/audit-logs", func(r chi.Router) {
                                r.With(authMiddleware.RequirePermission("audit_logs", "read")).Get("/", auditHandler.List)
                        })
                })
        })

        server := &http.Server{
                Addr:         ":" + cfg.Server.Port,
                Handler:      r,
                ReadTimeout:  cfg.Server.ReadTimeout,
                WriteTimeout: cfg.Server.WriteTimeout,
                IdleTimeout:  cfg.Server.IdleTimeout,
        }

        go func() {
                logger.Info().Str("port", cfg.Server.Port).Msg("Server is starting")
                if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                        logger.Fatal().Err(err).Msg("Server failed to start")
                }
        }()

        quit := make(chan os.Signal, 1)
        signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
        <-quit

        logger.Info().Msg("Server is shutting down...")

        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        if err := server.Shutdown(ctx); err != nil {
                logger.Fatal().Err(err).Msg("Server forced to shutdown")
        }

        logger.Info().Msg("Server exited properly")
}
