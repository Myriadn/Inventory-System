package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"project-app-inventory-restapi-golang-anas/config"
	"project-app-inventory-restapi-golang-anas/pkg/database"
	"project-app-inventory-restapi-golang-anas/pkg/logger"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func main() {
	// Load Config
	cfg := config.LoadConfig()

	// Setup Logger
	log := logger.NewLogger(cfg)
	defer log.Sync() // Flush buffer sebelum app mati

	// Connect Database
	db := database.ConnectDB(cfg)
	defer db.Close()

	// Setup Router (Chi)
	r := chi.NewRouter()

	// Middleware Dasar
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger) // Chi default logger (bisa diganti custom middleware zap nanti)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Test Route
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// TODO: Setup Routes/Handlers disini nanti

	// Start Server dengan Graceful Shutdown
	// Ini best practice agar request yang sedang berjalan tidak terputus paksa saat server dimatikan.
	server := &http.Server{
		Addr:    cfg.AppPort,
		Handler: r,
	}

	// Channel untuk listen sinyal interupsi (Ctrl+C atau kill command)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Jalankan server di goroutine terpisah
	go func() {
		log.Info("Starting server", zap.String("port", cfg.AppPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server startup failed", zap.Error(err))
		}
	}()

	// Block main thread sampai terima sinyal stop
	<-stop
	log.Info("Shutting down server...")

	// Context timeout untuk shutdown (maksimal 5 detik untuk selesaikan request sisa)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited properly")
}
