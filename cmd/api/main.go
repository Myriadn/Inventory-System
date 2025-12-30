package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"project-app-inventory-restapi-golang-anas/config"
	"project-app-inventory-restapi-golang-anas/internal/handler"
	"project-app-inventory-restapi-golang-anas/internal/middleware"
	"project-app-inventory-restapi-golang-anas/internal/repository"
	"project-app-inventory-restapi-golang-anas/internal/service"
	"project-app-inventory-restapi-golang-anas/pkg/database"
	"project-app-inventory-restapi-golang-anas/pkg/logger"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
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

	// Repository
	userRepo := repository.NewUserRepository(db)

	// Services
	authService := service.NewAuthService(userRepo)

	// Handler
	authHandler := handler.NewAuthHandler(authService)

	// Middleware auth
	authMiddleware := middleware.NewAuthMiddleware(userRepo)

	// Setup Router (Chi)
	r := chi.NewRouter()

	// Middleware Dasar
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger) // Chi default logger (bisa diganti custom middleware zap nanti)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// Test Route
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Setup Routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Group(func(r chi.Router) {
		// Pasang Middleware VerifyToken disini
		r.Use(authMiddleware.VerifyToken)

		// Contoh Endpoint User Profile (Bisa diakses semua role yang login)
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(middleware.UserIDKey).(int64)
			role := r.Context().Value(middleware.RoleKey).(string)
			w.Write([]byte(fmt.Sprintf("Hello User ID: %d, Role: %s", userID, role)))
		})

		// Admin Only Routes
		r.Group(func(r chi.Router) {
			// Pasang Middleware Role Check disini
			r.Use(authMiddleware.RequireRoles("super_admin", "admin"))

			r.Get("/admin-dashboard", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Welcome Admin!"))
			})
		})

		// Super Admin Only
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireRoles("super_admin"))
			// Disini nanti tempat CRUD User (Create/Delete Staff)
		})
	})

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
