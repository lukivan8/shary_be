package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"shary_be/internal/models"
	"shary_be/internal/service"

	"github.com/go-chi/chi/v5"

	"go.uber.org/zap"
)

type ItemPhotoHandler struct {
	itemPhotoService *service.ItemPhotoService
	logger           *zap.Logger
}

func NewItemPhotoHandler(itemPhotoService *service.ItemPhotoService, logger *zap.Logger) *ItemPhotoHandler {
	return &ItemPhotoHandler{
		itemPhotoService: itemPhotoService,
		logger:           logger,
	}
}

// GetPhotosByItemID handles GET /api/item_photos/{item_id}
func (h *ItemPhotoHandler) GetPhotosByItemID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract item ID from Chi URL parameters
	itemIDStr := chi.URLParam(r, "item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	photos, err := h.itemPhotoService.GetPhotosByItemID(itemID)
	if err != nil {
		h.logger.Error("Failed to get photos by item ID", zap.Int("item_id", itemID), zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"photos": photos,
	})
}

func (h *ItemPhotoHandler) AddPhotos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract item ID from Chi URL parameters
	itemIDStr := chi.URLParam(r, "item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var req models.CreateItemPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.itemPhotoService.AddPhotos(itemID, req.Photos); err != nil {
		h.logger.Error("Failed to add photos", zap.Int("item_id", itemID), zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"item_id": itemID,
	})
}

func (h *ItemPhotoHandler) DeletePhotos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract item ID from Chi URL parameters
	itemIDStr := chi.URLParam(r, "item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var req models.DeleteItemPhotosRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.itemPhotoService.DeletePhotos(itemID, req.PhotoIDs); err != nil {
		h.logger.Error("Failed to delete photos", zap.Int("item_id", itemID), zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ItemPhotoHandler) CountPhotosByItemID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract item ID from Chi URL parameters
	itemIDStr := chi.URLParam(r, "item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	count, err := h.itemPhotoService.CountPhotosByItemID(itemID)
	if err != nil {
		h.logger.Error("Failed to count photos", zap.Int("item_id", itemID), zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": count,
	})
}
