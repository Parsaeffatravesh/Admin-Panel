package config

import (
        "os"
        "strconv"
        "strings"
        "time"
)

type Config struct {
        Server   ServerConfig
        Database DatabaseConfig
        JWT      JWTConfig
        App      AppConfig
}

type ServerConfig struct {
        Port           string
        ReadTimeout    time.Duration
        WriteTimeout   time.Duration
        IdleTimeout    time.Duration
        AllowedOrigins []string
}

type DatabaseConfig struct {
        URL             string
        MaxConns        int32
        MinConns        int32
        MaxConnLifetime time.Duration
        MaxConnIdleTime time.Duration
}

type JWTConfig struct {
        Secret          string
        AccessTokenTTL  time.Duration
        RefreshTokenTTL time.Duration
}

type AppConfig struct {
        Environment string
        LogLevel    string
}

func Load() *Config {
        allowedOrigins := getStringSliceEnv("ALLOWED_ORIGINS", nil)
        if len(allowedOrigins) == 0 {
                allowedOrigins = []string{"*"}
        }
        
        return &Config{
                Server: ServerConfig{
                        Port:           getEnv("PORT", "8080"),
                        ReadTimeout:    getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
                        WriteTimeout:   getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
                        IdleTimeout:    getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
                        AllowedOrigins: allowedOrigins,
                },
                Database: DatabaseConfig{
                        URL:             getEnv("DATABASE_URL", ""),
                        MaxConns:        int32(getIntEnv("DB_MAX_CONNS", 50)),
                        MinConns:        int32(getIntEnv("DB_MIN_CONNS", 10)),
                        MaxConnLifetime: getDurationEnv("DB_MAX_CONN_LIFETIME", 30*time.Minute),
                        MaxConnIdleTime: getDurationEnv("DB_MAX_CONN_IDLE_TIME", 10*time.Minute),
                },
                JWT: JWTConfig{
                        Secret:          getEnv("SESSION_SECRET", "default-secret-change-me"),
                        AccessTokenTTL:  getDurationEnv("JWT_ACCESS_TTL", 15*time.Minute),
                        RefreshTokenTTL: getDurationEnv("JWT_REFRESH_TTL", 7*24*time.Hour),
                },
                App: AppConfig{
                        Environment: getEnv("APP_ENV", "development"),
                        LogLevel:    getEnv("LOG_LEVEL", "debug"),
                },
        }
}

func getEnv(key, defaultValue string) string {
        if value := os.Getenv(key); value != "" {
                return value
        }
        return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
        if value := os.Getenv(key); value != "" {
                if intVal, err := strconv.Atoi(value); err == nil {
                        return intVal
                }
        }
        return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
        if value := os.Getenv(key); value != "" {
                if duration, err := time.ParseDuration(value); err == nil {
                        return duration
                }
        }
        return defaultValue
}

func getStringSliceEnv(key string, defaultValue []string) []string {
        if value := os.Getenv(key); value != "" {
                parts := strings.Split(value, ",")
                result := make([]string, 0, len(parts))
                for _, part := range parts {
                        trimmed := strings.TrimSpace(part)
                        if trimmed != "" {
                                result = append(result, trimmed)
                        }
                }
                if len(result) > 0 {
                        return result
                }
        }
        return defaultValue
}
