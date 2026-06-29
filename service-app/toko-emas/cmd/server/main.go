// @title           Jason Jewelry API
// @version         1.0

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"toko-emas/internal/conf"
	"toko-emas/internal/data"
	"toko-emas/internal/server"
	"toko-emas/internal/service"

	"gopkg.in/yaml.v3"
)

func main() {
	confDir := flag.String("conf", "configs", "config directory")
	flag.Parse()

	// Load config
	cfg, err := loadConfig(*confDir + "/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect database
	db, err := data.NewDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// Repositories
	barangRepo := data.NewBarangRepo(db)
	userRepo := data.NewUserRepo(db)
	pembelianRepo := data.NewPembelianRepo(db)
	karatRepo := data.NewKaratRepo(db)
	bakiRepo := data.NewBakiRepo(db)
	penjualanRepo := data.NewPenjualanRepo(db)

	// Services
	barangSvc := service.NewBarangService(barangRepo)
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(userRepo, cfg)
	pembelianSvc := service.NewPembelianService(pembelianRepo)
	karatSvc := service.NewKaratService(karatRepo)
	bakiSvc := service.NewBakiService(bakiRepo)
	penjualanSvc := service.NewPenjualanService(penjualanRepo)

	// Router
	router := server.NewHTTPRouter(barangSvc, userSvc, authSvc, pembelianSvc, karatSvc, bakiSvc, penjualanSvc)

	addr := cfg.Server.HTTP.Addr
	fmt.Printf("🚀 Server running on http://%s\n", addr)
	fmt.Printf("📖 Swagger UI: http://%s/docs/index.html\n", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func loadConfig(path string) (*conf.Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg conf.Config
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}
