package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/amr0ny/followers-etl-service/config"
	application "github.com/amr0ny/followers-etl-service/internal/app"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	conf, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	m, err := migrate.New(fmt.Sprintf("file://%s", conf.MigrationsDir), conf.DSN)
	if err != nil {
		log.Fatalf("failed to init migrations: %v", err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("failed to up migrations: %v", err)
	}

	app, err := application.InitializeApp(ctx, conf)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	if err := app.StartCron(ctx, conf.CronSchedule); err != nil {
		app.Logger.Error("failed to start cron", "error", err)
		return
	}
	<-ctx.Done()
}
