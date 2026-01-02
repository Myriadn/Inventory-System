package handler

import (
	"encoding/json"
	"net/http"
	"project-app-inventory-restapi-golang-anas/internal/entity"
	"project-app-inventory-restapi-golang-anas/internal/service"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type RackHandler struct {
	service   *service.RackService
	validator *validator.Validate
}

func NewRackHandler(service *service.RackService) *RackHandler {
	return &RackHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *RackHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req entity.RackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cat, err := h.service.Create(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cat)
}

func (h *RackHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Ambil query param page & limit
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	racks, err := h.service.GetAll(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle jika kosong biar return array kosong [], bukan null
	if racks == nil {
		racks = []entity.Rack{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(racks)
}

func (h *RackHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	rak, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rak)
}

func (h *RackHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var req entity.RackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Update(r.Context(), id, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "rack updated successfully"})
}

func (h *RackHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "rack deleted successfully"})
}
