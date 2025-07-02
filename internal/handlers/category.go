package handlers

import (
	"encoding/json"
	"net/http"

	"shary_be/internal/models"
	"shary_be/internal/service"

	"go.uber.org/zap"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
	logger *zap.Logger
}

func NewCategoryHandler(categoryService *service.CategoryService, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		logger: logger,
	}
}

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		h.logger.Error("Failed to get all categories", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(categories)
}

// CreateCategory handles POST /api/categories
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	category, err := h.categoryService.CreateCategory(&req)
	if err != nil {
		h.logger.Error("Failed to create category", zap.Error(err))

		if err.Error() == "category already exists" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}