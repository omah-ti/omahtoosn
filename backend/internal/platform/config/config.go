package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppName          string
	AppEnv           string
	AppPort          string
	AppVersion       string
	DatabaseURL      string
	JWTSecret        string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	CookieDomain     string
	CookieSecure     bool
	CookieSameSite   string
	CORSAllowOrigins string
}

func Load() *Config {
	return &Config{
		AppName:          getString("APP_NAME", "to-osn-backend"),
		AppEnv:           getString("APP_ENV", "development"),
		AppPort:          getString("APP_PORT", "8081"),
		AppVersion:       getString("APP_VERSION", "dev"),
		DatabaseURL:      getString("DATABASE_URL", "postgres://postgres:postgres@localhost:5433/to_osn?sslmode=disable"),
		JWTSecret:        getString("JWT_SECRET", "change-this-secret"),
		AccessTokenTTL:   time.Duration(getInt("ACCESS_TOKEN_TTL_MINUTES", 15)) * time.Minute,
		RefreshTokenTTL:  time.Duration(getInt("REFRESH_TOKEN_TTL_HOURS", 168)) * time.Hour,
		CookieDomain:     getString("COOKIE_DOMAIN", ""),
		CookieSecure:     getBool("COOKIE_SECURE", false),
		CookieSameSite:   normalizeSameSite(getString("COOKIE_SAME_SITE", "Lax")),
		CORSAllowOrigins: getString("CORS_ALLOW_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000"),
	}
}

func getString(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func getInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func normalizeSameSite(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "strict":
		return "Strict"
	case "none":
		return "None"
	default:
		return "Lax"
	}
}
