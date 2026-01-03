package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project-app-inventory-restapi-golang-anas/config"
	"project-app-inventory-restapi-golang-anas/internal/handler"
	"project-app-inventory-restapi-golang-anas/internal/middleware"
	"project-app-inventory-restapi-golang-anas/internal/repository"
	"project-app-inventory-restapi-golang-anas/internal/routes" // Import package routes baru
	"project-app-inventory-restapi-golang-anas/internal/service"
	"project-app-inventory-restapi-golang-anas/pkg/database"
	"project-app-inventory-restapi-golang-anas/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Load Config & Logger
	cfg := config.LoadConfig()
	log := logger.NewLogger(cfg)
	defer log.Sync()

	// Connect Database
	db := database.ConnectDB(cfg)
	defer db.Close()

	// Init Layers (Repository -> Service -> Handler)

	// User & Auth
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	authMiddleware := middleware.NewAuthMiddleware(userRepo)

	// Category
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Warehouse
	warehouseRepo := repository.NewWarehouseRepository(db)
	warehuseService := service.NewWarehouseService(warehouseRepo)
	warehouseHandler := handler.NewWarehouseHandler(warehuseService)

	// Rack
	rackRepo := repository.NewRackRepository(db)
	rackService := service.NewRackService(rackRepo)
	rackHandler := handler.NewRackHandler(rackService)

	// Product
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// Sale
	saleRepo := repository.NewSaleRepository(db)
	saleService := service.NewSaleService(saleRepo)
	saleHandler := handler.NewSaleHandler(saleService)

	// Report
	reportRepo := repository.NewReportRepository(db)
	reportService := service.NewReportService(reportRepo)
	reportHandler := handler.NewReportHandler(reportService)

	// Setup Router (Panggil fungsi dari package routes)
	// Kita passing semua dependency yang dibutuhkan router struct
	router := routes.SetupRouter(routes.RouterDeps{
		AuthMiddleware:   authMiddleware,
		AuthHandler:      authHandler,
		UserHandler:      userHandler,
		CategoryHandler:  categoryHandler,
		WarehouseHandler: warehouseHandler,
		RackHandler:      rackHandler,
		ProductHandler:   productHandler,
		SaleHandler:      saleHandler,
		ReportHandler:    reportHandler,
	})

	// Start Server (Graceful Shutdown)
	server := &http.Server{
		Addr:    cfg.AppPort,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("Starting server", zap.String("port", cfg.AppPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server startup failed", zap.Error(err))
		}
	}()

	<-stop
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited properly")
}
