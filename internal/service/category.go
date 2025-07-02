package service

import (
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
func (s *CategoryService) GetAllCategories() ([]	models.Category, error) {
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
		Name:        req.Name,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		s.logger.Error("Failed to create category", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Category created successfully", zap.Int("category_id", category.ID))
	return category, nil
}