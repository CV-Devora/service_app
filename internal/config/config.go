package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	AppName        string
	HTTPAddr       string
	DatabaseDriver string
	DatabaseURL    string
	MigrateOnStart bool
	SwaggerEnabled bool
}

func Load(configPath string) Config {
	values := loadConfigFile(configPath)

	return Config{
		AppName:        getEnvFromMap(values, "APP_NAME", "service-app"),
		HTTPAddr:       getEnvFromMap(values, "HTTP_ADDR", ":8000"),
		DatabaseDriver: getEnvFromMap(values, "DATABASE_DRIVER", "pgx"),
		DatabaseURL:    getEnvFromMap(values, "DATABASE_URL", "postgres://postgres:Tambunan140705@localhost:5432/jason_jewelry?sslmode=disable"),
		MigrateOnStart: getEnvFromMap(values, "MIGRATE_ON_START", "true") == "true",
		SwaggerEnabled: getEnvFromMap(values, "SWAGGER_ENABLED", "true") == "true",
	}
}

func loadConfigFile(configPath string) map[string]string {
	values := map[string]string{}
	if configPath == "" {
		return values
	}

	info, err := os.Stat(configPath)
	if err != nil {
		return values
	}

	filePath := configPath
	if info.IsDir() {
		filePath = filepath.Join(configPath, "app.env")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return values
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		values[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	return values
}

func getEnvFromMap(values map[string]string, key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if value := values[key]; value != "" {
		return value
	}
	return fallback
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
