package app

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"service-app/internal/config"
	"service-app/internal/infrastructure/postgres"
	httptransport "service-app/internal/transport/http"
	"service-app/internal/usecase"
)

type Container struct {
	db     *sql.DB
	server *httptransport.Server
}

func NewContainer(cfg config.Config) (*Container, error) {
	db, err := sql.Open(cfg.DatabaseDriver, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	if cfg.MigrateOnStart {
		goose.SetDialect("postgres")
		if err := goose.Up(db, "migrations"); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("run migrations: %w", err)
		}
	}

	userRepo := postgres.NewUserRepository(db)
	userSvc := usecase.NewUserService(userRepo)
	userHandler := httptransport.NewUserHandler(userSvc)

	srv := httptransport.NewServer(cfg.HTTPAddr, cfg.SwaggerEnabled, userHandler)

	return &Container{
		db:     db,
		server: srv,
	}, nil
}

func (c *Container) Run(ctx context.Context) error {
	return c.server.Run(ctx)
}

func (c *Container) Stop(ctx context.Context) error {
	if err := c.server.Stop(ctx); err != nil {
		return err
	}
	return c.db.Close()
}
