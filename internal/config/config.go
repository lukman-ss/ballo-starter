package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds runtime configuration loaded from environment variables.
type Config struct {
	AppEnv                string
	HTTPAddr              string
	HTTPReadHeaderTimeout time.Duration
	HTTPReadTimeout       time.Duration
	HTTPWriteTimeout      time.Duration
	HTTPIdleTimeout       time.Duration
	ShutdownTimeout       time.Duration
	MaxBodyBytes          int64
	LogLevel              string
	// Database (optional)
	DatabaseDriver   string
	DatabaseDSN      string
	DatabaseMaxConns int
	DatabaseMaxIdle  int
}

// Load loads configuration from environment variables with safe defaults.
func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:                getEnv("APP_ENV", "development"),
		HTTPAddr:              getEnv("HTTP_ADDR", ":8080"),
		HTTPReadHeaderTimeout: 5 * time.Second,
		HTTPReadTimeout:       10 * time.Second,
		HTTPWriteTimeout:      10 * time.Second,
		HTTPIdleTimeout:       120 * time.Second,
		ShutdownTimeout:       10 * time.Second,
		MaxBodyBytes:          1048576, // 1 MB
		LogLevel:              getEnv("LOG_LEVEL", "info"),
		DatabaseDriver:        getEnv("DATABASE_DRIVER", ""),
		DatabaseDSN:           getEnv("DATABASE_DSN", ""),
	}

	if v := os.Getenv("HTTP_READ_HEADER_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTP_READ_HEADER_TIMEOUT %q: %w", v, err)
		}
		cfg.HTTPReadHeaderTimeout = d
	}

	if v := os.Getenv("HTTP_READ_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTP_READ_TIMEOUT %q: %w", v, err)
		}
		cfg.HTTPReadTimeout = d
	}

	if v := os.Getenv("HTTP_WRITE_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTP_WRITE_TIMEOUT %q: %w", v, err)
		}
		cfg.HTTPWriteTimeout = d
	}

	if v := os.Getenv("HTTP_IDLE_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTP_IDLE_TIMEOUT %q: %w", v, err)
		}
		cfg.HTTPIdleTimeout = d
	}

	if v := os.Getenv("SHUTDOWN_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid SHUTDOWN_TIMEOUT %q: %w", v, err)
		}
		cfg.ShutdownTimeout = d
	}

	if v := os.Getenv("DATABASE_MAX_CONNS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return nil, fmt.Errorf("invalid DATABASE_MAX_CONNS %q: must be non‑negative integer", v)
		}
		cfg.DatabaseMaxConns = n
	}

	if v := os.Getenv("DATABASE_MAX_IDLE"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return nil, fmt.Errorf("invalid DATABASE_MAX_IDLE %q: must non‑negative integer", v)
		}
		cfg.DatabaseMaxIdle = n
	}

	if v := os.Getenv("MAX_BODY_BYTES"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n <= 0 {
			return nil, fmt.Errorf("invalid MAX_BODY_BYTES %q: must be positive integer", v)
		}
		cfg.MaxBodyBytes = n
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks configuration sanity.
func (c *Config) Validate() error {
	switch strings.ToLower(c.AppEnv) {
	case "development", "production", "test", "staging":
	default:
		return fmt.Errorf("invalid APP_ENV %q: must be development, production, test, or staging", c.AppEnv)
	}

	if strings.TrimSpace(c.HTTPAddr) == "" {
		return fmt.Errorf("HTTP_ADDR cannot be empty")
	}

	if c.HTTPReadHeaderTimeout <= 0 {
		return fmt.Errorf("HTTP_READ_HEADER_TIMEOUT must be positive")
	}
	if c.HTTPReadTimeout <= 0 {
		return fmt.Errorf("HTTP_READ_TIMEOUT must be positive")
	}
	if c.HTTPWriteTimeout <= 0 {
		return fmt.Errorf("HTTP_WRITE_TIMEOUT must be positive")
	}
	if c.HTTPIdleTimeout <= 0 {
		return fmt.Errorf("HTTP_IDLE_TIMEOUT must be positive")
	}
	if c.ShutdownTimeout <= 0 {
		return fmt.Errorf("SHUTDOWN_TIMEOUT must be positive")
	}
	if c.DatabaseDriver != "" && c.DatabaseDSN == "" {
		return fmt.Errorf("DATABASE_DRIVER set but DATABASE_DSN empty")
	}
	if c.DatabaseDSN != "" && c.DatabaseDriver == "" {
		return fmt.Errorf("DATABASE_DSN set but DATABASE_DRIVER empty")
	}

	if c.MaxBodyBytes <= 0 {
		return fmt.Errorf("MAX_BODY_BYTES must be positive")
	}

	switch strings.ToLower(c.LogLevel) {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("invalid LOG_LEVEL %q: must be debug, info, warn, or error", c.LogLevel)
	}

	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
