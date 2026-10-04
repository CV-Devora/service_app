package conf

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is the top-level configuration structure
type Config struct {
	Server Server `yaml:"server"`
	Data   Data   `yaml:"data"`
	Log    Log    `yaml:"log"`
	Auth   Auth   `yaml:"auth"`
}

type Server struct {
	HTTP HTTP `yaml:"http"`
}

type HTTP struct {
	Addr    string `yaml:"addr"`
	Timeout string `yaml:"timeout"`
}

type Data struct {
	Database Database `yaml:"database"`
}

type Database struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type Log struct {
	Level string `yaml:"level"`
}

type Auth struct {
	JWTSecret      string `yaml:"jwt_secret"`
	AccessTokenTTL string `yaml:"access_token_ttl"`
}

// Load loads configuration from a YAML file (if present) and overrides fields
// with environment variables (e.g. Railway, Docker, or .env).
func Load(path string) (*Config, error) {
	// Attempt to load .env file if it exists (for local development)
	loadDotEnv(".env")

	cfg := &Config{
		Server: Server{
			HTTP: HTTP{
				Addr:    "0.0.0.0:8001",
				Timeout: "30s",
			},
		},
		Data: Data{
			Database: Database{
				Driver: "postgres",
			},
		},
		Log: Log{
			Level: "info",
		},
		Auth: Auth{
			JWTSecret:      "change-this-secret",
			AccessTokenTTL: "6h",
		},
	}

	// Resolve config file path (supports directory or direct file path)
	filePath := path
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		filePath = filepath.Join(path, "config.yaml")
	} else if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
		candidate := filepath.Join(path, "config.yaml")
		if _, err := os.Stat(candidate); err == nil {
			filePath = candidate
		}
	}

	// Read YAML config if file exists
	if f, err := os.ReadFile(filePath); err == nil {
		if err := yaml.Unmarshal(f, cfg); err != nil {
			return nil, fmt.Errorf("parse config file %s: %w", filePath, err)
		}
		log.Printf("[Config] Base configuration loaded from %s", filePath)
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("read config file %s: %w", filePath, err)
	} else {
		log.Printf("[Config] Config file not found at %s, using defaults and environment variables", filePath)
	}

	// Apply environment variable overrides (Railway / Cloud deployment)
	applyEnvOverrides(cfg)

	return cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	// ── Server HTTP overrides ────────────────────────────────────────────────
	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		cfg.Server.HTTP.Addr = "0.0.0.0:" + port
		log.Printf("[Config] Server address set from Railway PORT: %s", cfg.Server.HTTP.Addr)
	} else if addr := getEnvAny("SERVER_ADDR", "HTTP_ADDR"); addr != "" {
		cfg.Server.HTTP.Addr = addr
	}

	if timeout := getEnvAny("SERVER_TIMEOUT"); timeout != "" {
		cfg.Server.HTTP.Timeout = timeout
	}

	// ── Database overrides ───────────────────────────────────────────────────
	if driver := getEnvAny("DB_DRIVER", "DATABASE_DRIVER"); driver != "" {
		cfg.Data.Database.Driver = driver
	}

	// Priority 1: Full URL / DSN strings (standard Railway variables)
	dbURL := getEnvAny("DATABASE_URL", "DATABASE_PUBLIC_URL", "POSTGRES_URL", "DATABASE_DSN", "DB_DSN")
	if dbURL != "" {
		// If sslmode is explicitly specified in env and not in URL, append it
		sslmode := getEnvAny("PGSSLMODE", "DB_SSLMODE")
		if sslmode != "" && !strings.Contains(dbURL, "sslmode=") {
			sep := "?"
			if strings.Contains(dbURL, "?") {
				sep = "&"
			}
			dbURL = dbURL + sep + "sslmode=" + sslmode
		}
		cfg.Data.Database.DSN = dbURL
		log.Println("[Config] Database DSN configured from environment variable (DATABASE_URL / POSTGRES_URL)")
	} else if host := getEnvAny("PGHOST", "POSTGRES_HOST", "DB_HOST"); host != "" {
		// Priority 2: Individual variables provided by Railway PostgreSQL plugin
		port := getEnvAny("PGPORT", "POSTGRES_PORT", "DB_PORT")
		if port == "" {
			port = "5432"
		}
		user := getEnvAny("PGUSER", "POSTGRES_USER", "DB_USER")
		if user == "" {
			user = "postgres"
		}
		password := getEnvAny("PGPASSWORD", "POSTGRES_PASSWORD", "DB_PASSWORD")
		dbname := getEnvAny("PGDATABASE", "POSTGRES_DB", "DB_NAME")
		sslmode := getEnvAny("PGSSLMODE", "POSTGRES_SSLMODE", "DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}
		tz := getEnvAny("PGTZ", "POSTGRES_TZ", "DB_TIMEZONE", "TZ")
		if tz == "" {
			tz = "Asia/Jakarta"
		}

		cfg.Data.Database.DSN = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			host, user, password, dbname, port, sslmode, tz,
		)
		log.Printf("[Config] Database DSN constructed from individual environment variables (host=%s, port=%s, db=%s, user=%s)",
			host, port, dbname, user)
	}

	// ── Log overrides ────────────────────────────────────────────────────────
	if logLevel := getEnvAny("LOG_LEVEL"); logLevel != "" {
		cfg.Log.Level = logLevel
	}

	// ── Auth overrides ───────────────────────────────────────────────────────
	if jwtSecret := getEnvAny("JWT_SECRET"); jwtSecret != "" {
		cfg.Auth.JWTSecret = jwtSecret
	}
	if ttl := getEnvAny("ACCESS_TOKEN_TTL"); ttl != "" {
		cfg.Auth.AccessTokenTTL = ttl
	}
}

func getEnvAny(keys ...string) string {
	for _, key := range keys {
		if val := strings.TrimSpace(os.Getenv(key)); val != "" {
			return val
		}
	}
	return ""
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)
			if os.Getenv(key) == "" {
				os.Setenv(key, val)
			}
		}
	}
}
