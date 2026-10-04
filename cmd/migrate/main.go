package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"toko-emas/internal/conf"
	"toko-emas/internal/data"

	"gorm.io/gorm"
)

type Migration struct {
	Version string
	Name    string
	UpSQL   string
	DownSQL string
}

var gooseUpRe = regexp.MustCompile(`(?s)-- \+goose Up\r?\n-- \+goose StatementBegin\r?\n(.*?)-- \+goose StatementEnd`)
var gooseDownRe = regexp.MustCompile(`(?s)-- \+goose Down\r?\n-- \+goose StatementBegin\r?\n(.*?)-- \+goose StatementEnd`)

func main() {
	confDir := flag.String("conf", "configs", "config directory")
	cmd := flag.String("cmd", "up", "migration command: up, down, status, reset")
	flag.Parse()

	cfg, err := conf.Load(*confDir)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := data.NewDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	migrations := loadMigrations()

	switch *cmd {
	case "up":
		runUp(db, migrations)
	case "down":
		runDown(db, migrations, 1)
	case "reset":
		runDown(db, migrations, len(migrations))
	case "status":
		showStatus(db, migrations)
	default:
		log.Fatalf("unknown command: %s", *cmd)
	}
}

func loadMigrations() []Migration {
	entries, err := os.ReadDir("migrations")
	if err != nil {
		log.Fatalf("read migrations dir: %v", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	var migrations []Migration
	for _, f := range files {
		content, err := os.ReadFile(filepath.Join("migrations", f))
		if err != nil {
			log.Fatalf("read %s: %v", f, err)
		}

		parts := strings.SplitN(f, "_", 2)
		version := parts[0]
		name := strings.TrimSuffix(parts[1], ".sql")

		mig := Migration{
			Version: version,
			Name:    name,
		}

		if matches := gooseUpRe.FindStringSubmatch(string(content)); len(matches) > 1 {
			mig.UpSQL = strings.TrimSpace(matches[1])
		} else {
			log.Fatalf("no -- +goose Up section found in %s", f)
		}

		if matches := gooseDownRe.FindStringSubmatch(string(content)); len(matches) > 1 {
			mig.DownSQL = strings.TrimSpace(matches[1])
		}

		migrations = append(migrations, mig)
	}

	return migrations
}

func getAppliedVersions(db *gorm.DB) map[string]bool {
	var rows []struct{ Version string }
	db.Table("schema_migrations").Select("version").Find(&rows)
	applied := make(map[string]bool)
	for _, r := range rows {
		applied[r.Version] = true
	}
	return applied
}

func runUp(db *gorm.DB, migrations []Migration) {
	applied := getAppliedVersions(db)
	for _, m := range migrations {
		if applied[m.Version] {
			fmt.Printf("[SKIP] %s_%s (already applied)\n", m.Version, m.Name)
			continue
		}
		fmt.Printf("[UP]   %s_%s\n", m.Version, m.Name)
		if err := db.Exec(m.UpSQL).Error; err != nil {
			log.Fatalf("migration %s failed: %v", m.Version, err)
		}
		db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", m.Version)
	}
	fmt.Println("migrations applied successfully")
}

func runDown(db *gorm.DB, migrations []Migration, steps int) {
	applied := getAppliedVersions(db)
	count := 0
	for i := len(migrations) - 1; i >= 0 && count < steps; i-- {
		m := migrations[i]
		if !applied[m.Version] {
			continue
		}
		if m.DownSQL == "" {
			log.Fatalf("no down migration for %s_%s", m.Version, m.Name)
		}
		fmt.Printf("[DOWN] %s_%s\n", m.Version, m.Name)
		if err := db.Exec(m.DownSQL).Error; err != nil {
			log.Fatalf("migration down %s failed: %v", m.Version, err)
		}
		db.Exec("DELETE FROM schema_migrations WHERE version = ?", m.Version)
		count++
	}
	fmt.Println("rollback completed")
}

func showStatus(db *gorm.DB, migrations []Migration) {
	applied := getAppliedVersions(db)
	fmt.Println("Migration Status:")
	fmt.Println("=================")
	for _, m := range migrations {
		status := "PENDING"
		if applied[m.Version] {
			status = "APPLIED"
		}
		fmt.Printf("  %s_%-50s %s\n", m.Version, m.Name, status)
	}
}
