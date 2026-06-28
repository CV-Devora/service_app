package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"service-app/internal/app"
	"service-app/internal/config"
)

func main() {
	confPath := flag.String("conf", "", "path to config file or directory")
	flag.Parse()

	cfg := config.Load(*confPath)
	log.Printf("starting %s on %s", cfg.AppName, cfg.HTTPAddr)

	container, err := app.NewContainer(cfg)
	if err != nil {
		log.Fatalf("build app container: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := container.Run(ctx); err != nil {
			log.Printf("server stopped with error: %v", err)
			cancel()
		}
	}()

	log.Printf("server is running")
	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := container.Stop(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
