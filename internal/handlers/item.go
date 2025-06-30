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

// ItemHandler handles HTTP requests for items
type ItemHandler struct {
	itemService *service.ItemService
	logger      *zap.Logger
}

// NewItemHandler creates a new item handler
func NewItemHandler(itemService *service.ItemService, logger *zap.Logger) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
		logger:      logger,
	}
}

// GetAllItems handles GET /api/items with optional query parameters
func (h *ItemHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters for filtering
	filter := &models.ItemFilter{}

	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filter.MinPrice = &minPrice
		}
	}

	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filter.MaxPrice = &maxPrice
		}
	}

	if location := r.URL.Query().Get("location"); location != "" {
		filter.Location = &location
	}

	if availableStr := r.URL.Query().Get("available"); availableStr != "" {
		if available, err := strconv.ParseBool(availableStr); err == nil {
			filter.IsAvailable = &available
		}
	}

	if search := r.URL.Query().Get("search"); search != "" {
		filter.Search = &search
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

	items, err := h.itemService.GetAllItems(filter)
	if err != nil {
		h.logger.Error("Failed to get all items", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"items":   items,
		"count":   len(items),
		"filters": filter,
	})
}

// CreateItem handles POST /api/items
func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	item, err := h.itemService.CreateItem(&req)
	if err != nil {
		h.logger.Error("Failed to create item", zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// GetItemByID handles GET /api/items/{id}
func (h *ItemHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract item ID from Chi URL parameters
	itemIDStr := chi.URLParam(r, "id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	item, err := h.itemService.GetItemByID(itemID)
	if err != nil {
		h.logger.Error("Failed to get item by ID", zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(item)
}

// UpdateItem handles PUT /api/items/{id}
func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract item ID from Chi URL parameters
	itemIDStr := chi.URLParam(r, "id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	item, err := h.itemService.UpdateItem(itemID, &req)
	if err != nil {
		h.logger.Error("Failed to update item", zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(item)
}

// DeleteItem handles DELETE /api/items/{id}
func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	// Extract item ID from Chi URL parameters
	itemIDStr := chi.URLParam(r, "id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	err = h.itemService.DeleteItem(itemID)
	if err != nil {
		h.logger.Error("Failed to delete item", zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetItemsByLocation handles GET /api/items/location/{location}
func (h *ItemHandler) GetItemsByLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract location from Chi URL parameters
	location := chi.URLParam(r, "location")
	if location == "" {
		http.Error(w, "Location cannot be empty", http.StatusBadRequest)
		return
	}

	items, err := h.itemService.GetItemsByLocation(location)
	if err != nil {
		h.logger.Error("Failed to get items by location", zap.Error(err))

		if err.Error() == "location cannot be empty" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"items":    items,
		"count":    len(items),
		"location": location,
	})
}
