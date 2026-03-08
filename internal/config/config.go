package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHTTPAddr       = ":8080"
	defaultShutdownSecond = 10
)

// Config stores runtime settings for bms service.
type Config struct {
	HTTPAddr         string
	ShutdownTimeout  time.Duration
	AllowedOrigins   []string
	IncludeDocsRoute bool
	DBDSN            string
	JWTSecret        string
	AccessTTL        time.Duration
	RefreshTTL       time.Duration
}

func Load() Config {
	cfg := Config{
		HTTPAddr:         getEnv("BMS_HTTP_ADDR", defaultHTTPAddr),
		ShutdownTimeout:  time.Duration(getEnvInt("BMS_SHUTDOWN_TIMEOUT_SEC", defaultShutdownSecond)) * time.Second,
		AllowedOrigins:   parseCSV(getEnv("BMS_CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")),
		IncludeDocsRoute: getEnvBool("BMS_ENABLE_DOCS", true),
		DBDSN:            strings.TrimSpace(getEnv("BMS_DB_DSN", "file:bms.db?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)")),
		JWTSecret:        getEnv("BMS_JWT_SECRET", "dev-secret-change-me"),
		AccessTTL:        time.Duration(getEnvInt("BMS_ACCESS_TTL_SEC", 900)) * time.Second,
		RefreshTTL:       time.Duration(getEnvInt("BMS_REFRESH_TTL_SEC", 2592000)) * time.Second,
	}
	if cfg.ShutdownTimeout <= 0 {
		cfg.ShutdownTimeout = defaultShutdownSecond * time.Second
	}
	if cfg.AccessTTL <= 0 {
		cfg.AccessTTL = 15 * time.Minute
	}
	if cfg.RefreshTTL <= 0 {
		cfg.RefreshTTL = 30 * 24 * time.Hour
	}
	return cfg
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		v = strings.TrimSpace(v)
		if v != "" {
			return v
		}
	}
	return def
}

func getEnvInt(key string, def int) int {
	raw := getEnv(key, "")
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	return v
}

func getEnvBool(key string, def bool) bool {
	raw := strings.ToLower(strings.TrimSpace(getEnv(key, "")))
	if raw == "" {
		return def
	}
	switch raw {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return def
	}
}

func parseCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	return out
}
