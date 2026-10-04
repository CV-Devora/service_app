package conf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear any env vars that might interfere
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("PGHOST")
	os.Unsetenv("PORT")
	os.Unsetenv("JWT_SECRET")

	cfg, err := Load("non_existent_dir")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Server.HTTP.Addr != "0.0.0.0:8001" {
		t.Errorf("expected 0.0.0.0:8001, got %s", cfg.Server.HTTP.Addr)
	}
	if cfg.Data.Database.Driver != "postgres" {
		t.Errorf("expected postgres driver, got %s", cfg.Data.Database.Driver)
	}
	if cfg.Auth.JWTSecret != "change-this-secret" {
		t.Errorf("expected default jwt secret, got %s", cfg.Auth.JWTSecret)
	}
}

func TestLoad_FromYAML(t *testing.T) {
	tmpDir := t.TempDir()
	yamlContent := `server:
  http:
    addr: 127.0.0.1:9000
data:
  database:
    driver: postgres
    dsn: "host=localhost user=test pass=test"
auth:
  jwt_secret: "my-yaml-secret"
`
	err := os.WriteFile(filepath.Join(tmpDir, "config.yaml"), []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to write yaml: %v", err)
	}

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Server.HTTP.Addr != "127.0.0.1:9000" {
		t.Errorf("expected 127.0.0.1:9000, got %s", cfg.Server.HTTP.Addr)
	}
	if cfg.Data.Database.DSN != "host=localhost user=test pass=test" {
		t.Errorf("expected custom dsn, got %s", cfg.Data.Database.DSN)
	}
	if cfg.Auth.JWTSecret != "my-yaml-secret" {
		t.Errorf("expected my-yaml-secret, got %s", cfg.Auth.JWTSecret)
	}
}

func TestLoad_RailwayDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://postgres:secretpassword@containers-us-west-123.railway.app:6543/railway")
	t.Setenv("PORT", "3000")
	t.Setenv("JWT_SECRET", "railway-jwt-secret")

	cfg, err := Load("configs")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedDSN := "postgresql://postgres:secretpassword@containers-us-west-123.railway.app:6543/railway"
	if cfg.Data.Database.DSN != expectedDSN {
		t.Errorf("expected %s, got %s", expectedDSN, cfg.Data.Database.DSN)
	}
	if cfg.Server.HTTP.Addr != "0.0.0.0:3000" {
		t.Errorf("expected 0.0.0.0:3000, got %s", cfg.Server.HTTP.Addr)
	}
	if cfg.Auth.JWTSecret != "railway-jwt-secret" {
		t.Errorf("expected railway-jwt-secret, got %s", cfg.Auth.JWTSecret)
	}
}

func TestLoad_RailwayIndividualEnv(t *testing.T) {
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("DATABASE_PUBLIC_URL")
	os.Unsetenv("POSTGRES_URL")
	t.Setenv("PGHOST", "postgres.railway.internal")
	t.Setenv("PGPORT", "5432")
	t.Setenv("PGUSER", "postgres")
	t.Setenv("PGPASSWORD", "mypassword")
	t.Setenv("PGDATABASE", "railway")
	t.Setenv("PGSSLMODE", "require")

	cfg, err := Load("configs")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedDSN := "host=postgres.railway.internal user=postgres password=mypassword dbname=railway port=5432 sslmode=require TimeZone=Asia/Jakarta"
	if cfg.Data.Database.DSN != expectedDSN {
		t.Errorf("expected %s, got %s", expectedDSN, cfg.Data.Database.DSN)
	}
}
