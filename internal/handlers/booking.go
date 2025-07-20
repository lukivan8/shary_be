package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"shary_be/internal/models"
	"shary_be/internal/service"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type RentHandler struct {
	rentService *service.RentService
	logger         *zap.Logger
}

func NewRentHandler(rentService *service.RentService, logger *zap.Logger) *RentHandler {
	return &RentHandler{
		rentService: rentService,
		logger:         logger,
	}
}

// GetAllRents handles GET /api/rents
func (h *RentHandler) GetAllRents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	filter := &models.RentFilter{}

	if minStartDateStr := r.URL.Query().Get("min_start_date"); minStartDateStr != "" {
		if minStartDate, err := time.Parse(time.RFC3339, minStartDateStr); err == nil {
			filter.MinStartDate = &minStartDate
		}
	}
	if maxStartDateStr := r.URL.Query().Get("max_start_date"); maxStartDateStr != "" {
		if maxStartDate, err := time.Parse(time.RFC3339, maxStartDateStr); err == nil {
			filter.MaxStartDate = &maxStartDate
		}
	}

	if minEndDateStr := r.URL.Query().Get("min_end_date"); minEndDateStr != "" {
		if minEndDate, err := time.Parse(time.RFC3339, minEndDateStr); err == nil {
			filter.MinEndDate = &minEndDate
		}
	}
	if maxEndDateStr := r.URL.Query().Get("max_end_date"); maxEndDateStr != "" {
		if maxEndDate, err := time.Parse(time.RFC3339, maxEndDateStr); err == nil {
			filter.MaxEndDate = &maxEndDate
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	rents, err := h.rentService.GetAllRents(filter)
	if err != nil {
		h.logger.Error("failed to get all rents", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rents": rents,
		"count":    len(rents),
		"filters":  filter,
	})
}	

// CreateRent handles POST /api/rents
func (h *RentHandler) CreateRent(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	rent, err := h.rentService.CreateRent(&req)
	if err != nil {
		h.logger.Error("failed to create rent", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rent)
}

// GetRentByItemID handles GET /api/rents/item/{item_id}
func (h *RentHandler) GetRentByItemID(w http.ResponseWriter, r *http.Request) {
	itemIDStr := chi.URLParam(r, "item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "invalid item ID", http.StatusBadRequest)
		return
	}

	rents, err := h.rentService.GetByItemID(itemID)
	if err != nil {
		h.logger.Error("failed to get rents by item ID", zap.Int("item_id", itemID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rents)
}

// GetRentByUserID handles GET /api/rents/user/{user_id}
func (h *RentHandler) GetRentByUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	rents, err := h.rentService.GetByUserID(userID)
	if err != nil {
		h.logger.Error("failed to get rents by user ID", zap.Int("user_id", userID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rents)
}	

// GetRentByID handles GET /api/rents/{id}
func (h *RentHandler) GetRentByID(w http.ResponseWriter, r *http.Request) {
	rentIDStr := chi.URLParam(r, "id")
	rentID, err := strconv.Atoi(rentIDStr)
	if err != nil {
		http.Error(w, "invalid rent ID", http.StatusBadRequest)
		return
	}

	rent, err := h.rentService.GetByID(rentID)
	if err != nil {
		h.logger.Error("failed to get rent by ID", zap.Int("rent_id", rentID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rent)
}

// UpdateRent handles PUT /api/rents/{id}
func (h *RentHandler) UpdateRent(w http.ResponseWriter, r *http.Request) {
	rentIDStr := chi.URLParam(r, "id")
	rentID, err := strconv.Atoi(rentIDStr)
	if err != nil {
		http.Error(w, "invalid rent ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateRentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	rent, err := h.rentService.Update(rentID, &req)
	if err != nil {
		h.logger.Error("failed to update rent", zap.Int("rent_id", rentID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rent)
}

// DeleteRent handles DELETE /api/rents/{id}
func (h *RentHandler) DeleteRent(w http.ResponseWriter, r *http.Request) {
	rentIDStr := chi.URLParam(r, "id")
	rentID, err := strconv.Atoi(rentIDStr)
	if err != nil {
		http.Error(w, "invalid rent ID", http.StatusBadRequest)
		return
	}

	err = h.rentService.Delete(rentID)
	if err != nil {
		h.logger.Error("failed to delete rent", zap.Int("rent_id", rentID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAvailableRents handles GET /api/rents/available
func (h *RentHandler) GetAvailableRents(w http.ResponseWriter, r *http.Request) {
	rents, err := h.rentService.GetAvailableRents()
	if err != nil {
		h.logger.Error("failed to get available rents", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rents)
}