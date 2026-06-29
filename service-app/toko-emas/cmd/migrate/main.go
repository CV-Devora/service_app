package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Data struct {
		Database struct {
			DSN string `yaml:"dsn"`
		} `yaml:"database"`
	} `yaml:"data"`
}

func main() {
	confDir := flag.String("conf", "configs", "config directory")
	command := flag.String("cmd", "up", "goose command: up, down, status, reset")
	flag.Parse()

	f, err := os.ReadFile(*confDir + "/config.yaml")
	if err != nil {
		log.Fatalf("read config: %v", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		log.Fatalf("parse config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Data.Database.DSN)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("set dialect: %v", err)
	}

	migrationsDir := "migrations"
	switch *command {
	case "up":
		err = goose.Up(db, migrationsDir)
	case "down":
		err = goose.Down(db, migrationsDir)
	case "status":
		err = goose.Status(db, migrationsDir)
	case "reset":
		err = goose.Reset(db, migrationsDir)
	default:
		fmt.Printf("Unknown command: %s\n", *command)
		os.Exit(1)
	}

	if err != nil {
		log.Fatalf("goose %s: %v", *command, err)
	}
	fmt.Printf("Migration '%s' completed successfully\n", *command)
}
