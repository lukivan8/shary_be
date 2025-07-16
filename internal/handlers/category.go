package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"shary_be/internal/models"
	"shary_be/internal/service"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
	logger          *zap.Logger
}

func NewCategoryHandler(categoryService *service.CategoryService, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		logger:          logger,
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

// UpdateCategory handles PUT /api/categories/:id
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categoryIDStr := chi.URLParam(r, "id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	category, err := h.categoryService.UpdateCategory(categoryID, &req)
	if err != nil {
		h.logger.Error("Failed to update category", zap.Error(err))

		if err.Error() == "category not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(category)
}

// DeleteCategory handles DELETE /api/categories/{id}
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := chi.URLParam(r, "id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	err = h.categoryService.DeleteCategory(categoryID)
	if err != nil {
		h.logger.Error("Failed to delete category", zap.Error(err))

		if err.Error() == "category not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCategoryByID handles GET /api/categories/{id}
func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categoryIDStr := chi.URLParam(r, "id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := h.categoryService.GetCategoryByID(categoryID)
	if err != nil {
		h.logger.Error("Failed to get category by ID", zap.Error(err))

		if err.Error() == "category not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(category)
}
