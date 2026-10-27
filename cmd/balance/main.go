package main

import (
	"context"
	"errors"
	"flag"
	"github.com/ekuzm/SDLC/internal/balance"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var version = "dev"

func main() {
	restore := flag.String("restore", "", "Восстановить из резервной копии при остановленном сервере")
	backup := flag.Bool("backup", false, "Создать резервную копию при остановленном сервере")
	flag.Parse()
	path := env("BALANCE_DATA", "data/balance.json")
	if *restore != "" {
		if e := balance.Restore(path, *restore); e != nil {
			slog.Error("restore failed", "error", e)
			os.Exit(1)
		}
		slog.Info("restored")
		return
	}
	store, e := balance.Open(path)
	if e != nil {
		slog.Error("storage open failed", "error", e)
		os.Exit(1)
	}
	defer store.Close()
	if e = store.Backup(time.Now()); e != nil {
		slog.Error("backup failed", "error", e)
		os.Exit(1)
	}
	if *backup {
		return
	}
	loc, e := time.LoadLocation("Europe/Minsk")
	if e != nil {
		slog.Error("timezone unavailable", "error", e)
		os.Exit(1)
	}
	cert, key := os.Getenv("BALANCE_TLS_CERT"), os.Getenv("BALANCE_TLS_KEY")
	if (cert == "") != (key == "") {
		slog.Error("both BALANCE_TLS_CERT and BALANCE_TLS_KEY are required")
		os.Exit(1)
	}
	now := func() time.Time { return time.Now().In(loc) }
	app, e := balance.NewServer(store, balance.Config{Secure: env("BALANCE_COOKIE_SECURE", "false") == "true" || cert != "", Now: now})
	if e != nil {
		slog.Error("application setup failed", "error", e)
		os.Exit(1)
	}
	server := &http.Server{Addr: env("BALANCE_ADDR", "127.0.0.1:8080"), Handler: app, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if e := store.Backup(now()); e != nil {
					slog.Error("daily backup failed", "error", e)
				}
			}
		}
	}()
	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdown)
	}()
	slog.Info("Баланс запущен", "address", server.Addr, "version", version)
	if cert != "" && key != "" {
		e = server.ListenAndServeTLS(cert, key)
	} else {
		e = server.ListenAndServe()
	}
	if e != nil && !errors.Is(e, http.ErrServerClosed) {
		slog.Error("server failed", "error", e)
		os.Exit(1)
	}
	if errors.Is(e, http.ErrServerClosed) {
		<-shutdownDone
	}
}
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
