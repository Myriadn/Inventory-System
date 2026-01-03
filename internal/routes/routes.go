package routes

import (
	"fmt"
	"net/http"
	"project-app-inventory-restapi-golang-anas/internal/handler"
	"project-app-inventory-restapi-golang-anas/internal/middleware"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type RouterDeps struct {
	AuthMiddleware   *middleware.AuthMiddleware
	AuthHandler      *handler.AuthHandler
	UserHandler      *handler.UserHandler
	CategoryHandler  *handler.CategoryHandler
	WarehouseHandler *handler.WarehouseHandler
	RackHandler      *handler.RackHandler
	ProductHandler   *handler.ProductHandler
	SaleHandler      *handler.SaleHandler
	ReportHandler    *handler.ReportHandler
}

func SetupRouter(deps RouterDeps) *chi.Mux {
	r := chi.NewRouter()

	// Middleware Dasar (Global)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// Test Route
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// --- Public Routes (Auth) ---
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", deps.AuthHandler.Register)
		r.Post("/login", deps.AuthHandler.Login)
		r.Post("/logout", deps.AuthHandler.Logout)
	})

	// --- Protected Routes (Butuh Login) ---
	r.Group(func(r chi.Router) {
		// Middleware Cek Token (Validasi Login)
		r.Use(deps.AuthMiddleware.VerifyToken)

		// Profile User (Semua Role Boleh Akses)
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(middleware.UserIDKey).(int64)
			role := r.Context().Value(middleware.RoleKey).(string)
			w.Write([]byte(fmt.Sprintf("Hello User ID: %d, Role: %s", userID, role)))
		})

		// Master Data (Category, Warehouse, Rack, Product)
		// Rule:
		// - Staff: Boleh READ Only.
		// - Admin/Super Admin: Boleh CRUD.

		// Categories
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", deps.CategoryHandler.GetAll)
			r.Get("/{id}", deps.CategoryHandler.GetByID)

			// Group Khusus Admin & Super Admin (Write Access)
			r.Group(func(r chi.Router) {
				r.Use(deps.AuthMiddleware.RequireRoles("super_admin", "admin"))
				r.Post("/", deps.CategoryHandler.Create)
				r.Put("/{id}", deps.CategoryHandler.Update)
				r.Delete("/{id}", deps.CategoryHandler.Delete)
			})
		})

		// Warehouses
		r.Route("/warehouses", func(r chi.Router) {
			r.Get("/", deps.WarehouseHandler.GetAll)
			r.Get("/{id}", deps.WarehouseHandler.GetByID)

			r.Group(func(r chi.Router) {
				r.Use(deps.AuthMiddleware.RequireRoles("super_admin", "admin"))
				r.Post("/", deps.WarehouseHandler.Create)
				r.Put("/{id}", deps.WarehouseHandler.Update)
				r.Delete("/{id}", deps.WarehouseHandler.Delete)
			})
		})

		// Racks
		r.Route("/racks", func(r chi.Router) {
			r.Get("/", deps.RackHandler.GetAll)
			r.Get("/{id}", deps.RackHandler.GetByID)

			r.Group(func(r chi.Router) {
				r.Use(deps.AuthMiddleware.RequireRoles("super_admin", "admin"))
				r.Post("/", deps.RackHandler.Create)
				r.Put("/{id}", deps.RackHandler.Update)
				r.Delete("/{id}", deps.RackHandler.Delete)
			})
		})

		// Products
		r.Route("/products", func(r chi.Router) {
			// Endpoint Khusus: Cek Stok Minimum (Staff Boleh Akses)
			r.With(deps.AuthMiddleware.RequireRoles("super_admin", "admin", "staff")).
				Get("/low-stock", deps.ProductHandler.GetLowStock)

			r.Get("/", deps.ProductHandler.GetAll)
			r.Get("/{id}", deps.ProductHandler.GetByID)

			// Admin Only: Create/Update/Delete
			r.Group(func(r chi.Router) {
				r.Use(deps.AuthMiddleware.RequireRoles("super_admin", "admin"))
				r.Post("/", deps.ProductHandler.Create)
				r.Put("/{id}", deps.ProductHandler.Update)
				r.Delete("/{id}", deps.ProductHandler.Delete)
			})
		})

		// Transactions (Sales)
		// Rule: Staff, Admin, Super Admin boleh Create Sale
		r.Route("/sales", func(r chi.Router) {
			r.Use(deps.AuthMiddleware.RequireRoles("super_admin", "admin", "staff"))

			r.Post("/", deps.SaleHandler.Create)
			r.Get("/", deps.SaleHandler.GetAll)
			r.Get("/{id}", deps.SaleHandler.GetByID)
		})

		// Reports (Revenue/Dashboard)
		// Rule: Hanya Admin & Super Admin (Staff dilarang lihat Revenue)
		r.Route("/reports", func(r chi.Router) {
			r.Use(deps.AuthMiddleware.RequireRoles("super_admin", "admin"))

			r.Get("/dashboard", deps.ReportHandler.GetDashboard)
		})

		// User Management
		// Rule: Hanya Super Admin
		r.Route("/users", func(r chi.Router) {
			r.Use(deps.AuthMiddleware.RequireRoles("super_admin"))

			r.Get("/", deps.UserHandler.GetAll)
			r.Patch("/{id}/role", deps.UserHandler.UpdateRole)
			r.Delete("/{id}", deps.UserHandler.Delete)
		})
	})

	return r
}
