package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/app"
	"github.com/pafthang/bms/internal/config"
	"github.com/pafthang/bms/internal/db"
)

func main() {
	cfg := config.Load()
	database, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer database.Close()

	engine := app.BuildEngine(cfg, database, app.BuildOptions{
		IncludeSystemRoutes: true,
		IncludeHealthRoutes: true,
	})

	srv := arc.NewServer(cfg.HTTPAddr, engine)
	errCh := make(chan error, 1)
	go func() {
		log.Printf("bms api starting on %s", cfg.HTTPAddr)
		errCh <- srv.Start()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("shutdown signal received: %s", sig)
	case err = <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server stopped with error: %v", err)
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
	log.Printf("server stopped")
}
