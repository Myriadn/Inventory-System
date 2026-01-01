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

	// User Settings
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(userRepo)

	// Category
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Warehouse
	warehouseRepo := repository.NewWarehouseRepository(db)
	warehuseService := service.NewWarehouseService(warehouseRepo)
	warehouseHandler := handler.NewWarehouseHandler(warehuseService)

	// Racks
	rackRepo := repository.NewRackRepository(db)
	rackService := service.NewRackService(rackRepo)
	rackHandler := handler.NewRackHandler(rackService)

	// product
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// sale
	saleRepo := repository.NewSaleRepository(db)
	saleService := service.NewSaleService(saleRepo)
	saleHandler := handler.NewSaleHandler(saleService)

	// report
	reportRepo := repository.NewReportRepository(db)
	reportService := service.NewReportService(reportRepo)
	reportHandler := handler.NewReportHandler(reportService)

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

		// categories
		r.Route("/categories", func(r chi.Router) {
			// Siapa saja (Admin, Staff) boleh READ
			r.Get("/", categoryHandler.GetAll)
			r.Get("/{id}", categoryHandler.GetByID)
			// Khusus ADMIN & SUPER ADMIN boleh CREATE/UPDATE/DELETE
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireRoles("super_admin", "admin"))
				r.Post("/", categoryHandler.Create)
				r.Put("/{id}", categoryHandler.Update)
				r.Delete("/{id}", categoryHandler.Delete)
			})
		})

		// warehouse
		r.Route("/warehouses", func(r chi.Router) {
			// Siapa saja (Admin, Staff) boleh READ
			r.Get("/", warehouseHandler.GetAll)
			r.Get("/{id}", warehouseHandler.GetByID)
			// Khusus ADMIN & SUPER ADMIN boleh CREATE/UPDATE/DELETE
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireRoles("super_admin", "admin"))
				r.Post("/", warehouseHandler.Create)
				r.Put("/{id}", warehouseHandler.Update)
				r.Delete("/{id}", warehouseHandler.Delete)
			})
		})

		// racks
		r.Route("/racks", func(r chi.Router) {
			// Siapa saja (Admin, Staff) boleh READ
			r.Get("/", rackHandler.GetAll)
			r.Get("/{id}", rackHandler.GetByID)
			// Khusus ADMIN & SUPER ADMIN boleh CREATE/UPDATE/DELETE
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireRoles("super_admin", "admin"))
				r.Post("/", rackHandler.Create)
				r.Put("/{id}", rackHandler.Update)
				r.Delete("/{id}", rackHandler.Delete)
			})
		})

		// products
		r.Route("/products", func(r chi.Router) {
			// Endpoint khusus Low Stock (Boleh Staff, Admin, SuperAdmin)
			r.With(authMiddleware.RequireRoles("super_admin", "admin", "staff")).
				Get("/low-stock", productHandler.GetLowStock)

			// Semua role boleh READ
			r.Get("/", productHandler.GetAll)
			r.Get("/{id}", productHandler.GetByID)

			// Hanya Admin & Super Admin boleh Create/Update/Delete
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireRoles("super_admin", "admin"))
				r.Post("/", productHandler.Create)
				r.Put("/{id}", productHandler.Update)
				r.Delete("/{id}", productHandler.Delete)
			})
		})

		// sale
		r.Route("/sales", func(r chi.Router) {
			// Admin & Staff boleh create sale
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireRoles("super_admin", "admin", "staff"))

				r.Post("/", saleHandler.Create)
			})
		})

		// reports (Hanya Admin & Super Admin)
		r.Route("/reports", func(r chi.Router) {
			r.Use(authMiddleware.RequireRoles("super_admin", "admin"))

			r.Get("/dashboard", reportHandler.GetDashboard)
		})

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
