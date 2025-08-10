package service

import (
	"database/sql"
	"shary_be/internal/models"
	"shary_be/internal/repository"

	"go.uber.org/zap"
)

// CategoryService handles business logic for categories
type CategoryService struct {
	categoryRepo *repository.CategoryRepository
	logger       *zap.Logger
}

func NewCategoryService(categoryRepo *repository.CategoryRepository, logger *zap.Logger) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

// GetAll retrieves all categories
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	categories, err := s.categoryRepo.GetAll()
	if err != nil {
		s.logger.Error("Failed to get all categories", zap.Error(err))
		return nil, err
	}

	return categories, nil
}

// CreateCategory adds a new category
func (s *CategoryService) CreateCategory(req *models.CreateCategoryRequest) (*models.Category, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		s.logger.Error("Invalid create category request", zap.Error(err))
		return nil, err
	}

	// Create category
	category := &models.Category{
		Name: req.Name,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		s.logger.Error("Failed to create category", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Category created successfully", zap.Int("category_id", category.ID))
	return category, nil
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(id int, req *models.UpdateCategoryRequest) (*models.Category, error) {
	if err := req.Validate(); err != nil {
		s.logger.Error("Invalid update category request", zap.Error(err))
		return nil, err
	}

	currentCategory, err := s.categoryRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get category for update", zap.Int("category_id", id), zap.Error(err))
		return nil, err
	}
	if currentCategory == nil {
		return nil, sql.ErrNoRows
	}

	categoryToUpdate := &models.Category{
		ID:   currentCategory.ID,
		Name: currentCategory.Name,
	}

	if req.Name != nil {
		categoryToUpdate.Name = *req.Name
	}

	if err := s.categoryRepo.Update(categoryToUpdate); err != nil {
		s.logger.Error("Failed to update category", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Category updated successfully", zap.Int("category_id", id))
	return categoryToUpdate, nil
}

// DeleteCategory deletes a category by ID
func (s *CategoryService) DeleteCategory(id int) error {
	// Delete category
	if err := s.categoryRepo.Delete(id); err != nil {
		s.logger.Error("Failed to delete category", zap.Error(err))
		return err
	}

	s.logger.Info("Category deleted successfully", zap.Int("category_id", id))
	return nil
}

// GetCategoryByID retrieves a category by ID
func (s *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get category by ID", zap.Int("category_id", id), zap.Error(err))
		return nil, err
	}

	if category == nil {
		return nil, sql.ErrNoRows
	}

	return category, nil
}
