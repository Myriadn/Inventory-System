package handler

import (
	"encoding/json"
	"net/http"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"project-app-inventory-restapi-golang-anas/internal/middleware"
	"project-app-inventory-restapi-golang-anas/internal/service"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type SaleHandler struct {
	service   *service.SaleService
	validator *validator.Validate
}

func NewSaleHandler(service *service.SaleService) *SaleHandler {
	return &SaleHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *SaleHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Ambil UserID dari Context (Hasil Login)
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse Request
	var req entity.SaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Panggil Service
	sale, err := h.service.CreateSale(r.Context(), userID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sale)
}

func (h *SaleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	sales, err := h.service.GetAll(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if sales == nil {
		sales = []entity.Sale{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sales)
}

func (h *SaleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	sale, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sale)
}
