// Package app wires the application using go-kratos bootstrap pattern.
package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"

	"toko-emas/internal/conf"
	"toko-emas/internal/data"
	"toko-emas/internal/server"
	"toko-emas/internal/service"

	"gopkg.in/yaml.v3"
)

// Run bootstraps and starts the application.
func Run(confDir string) {
	cfg, err := loadConfig(confDir + "/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// ── Data layer ──────────────────────────────────────────────────────────
	db, err := data.NewDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	barangRepo    := data.NewBarangRepo(db)
	userRepo      := data.NewUserRepo(db)
	pembelianRepo := data.NewPembelianRepo(db)
	karatRepo     := data.NewKaratRepo(db)
	bakiRepo      := data.NewBakiRepo(db)
	penjualanRepo := data.NewPenjualanRepo(db)
	dashboardRepo := data.NewDashboardRepo(db)
	barangLandingRepo := data.NewBarangLandingRepo(db)

	// ── Service layer ────────────────────────────────────────────────────────
	authSvc      := service.NewAuthService(userRepo, cfg)
	barangSvc    := service.NewBarangService(barangRepo, authSvc)
	userSvc      := service.NewUserService(userRepo)
	pembelianSvc := service.NewPembelianService(pembelianRepo)
	karatSvc     := service.NewKaratService(karatRepo)
	bakiSvc      := service.NewBakiService(bakiRepo)
	penjualanSvc := service.NewPenjualanService(penjualanRepo, authSvc)
	dashboardSvc := service.NewDashboardService(dashboardRepo)
	uploadSvc := service.NewUploadService()
	barangLandingSvc := service.NewBarangLandingService(barangLandingRepo)

	// ensure uploads directory exists
	os.MkdirAll("uploads", 0755)

	// ── HTTP router (gorilla/mux with Swagger) ───────────────────────────────
	handler := server.NewHTTPRouter(
		barangSvc, userSvc, authSvc,
		pembelianSvc, karatSvc, bakiSvc,
		penjualanSvc, dashboardSvc, uploadSvc,
		barangLandingSvc,
	)

	// ── go-kratos HTTP server ─────────────────────────────────────────────────
	hs := kratoshttp.NewServer(
		kratoshttp.Address(cfg.Server.HTTP.Addr),
		kratoshttp.Timeout(30*time.Second),
	)
	hs.HandlePrefix("/", handler)

	// ── Graceful lifecycle ───────────────────────────────────────────────────
	addr := cfg.Server.HTTP.Addr
	fmt.Printf("🚀  Toko Emas API    → http://%s\n", addr)
	fmt.Printf("📖  Swagger UI       → http://%s/docs/index.html\n", addr)

	if err := hs.Start(context.Background()); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server start error: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := hs.Stop(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("server stopped.")
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
