package service

import (
	"database/sql"
	"errors"

	"shary_be/internal/models"
	"shary_be/internal/repository"

	"go.uber.org/zap"
)

// ItemService handles business logic for items
type ItemService struct {
	itemRepo *repository.ItemRepository
	logger   *zap.Logger
}

// NewItemService creates a new item service
func NewItemService(itemRepo *repository.ItemRepository, logger *zap.Logger) *ItemService {
	return &ItemService{
		itemRepo: itemRepo,
		logger:   logger,
	}
}

// CreateItem creates a new item
func (s *ItemService) CreateItem(req *models.CreateItemRequest) (*models.Item, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		s.logger.Error("Invalid create item request", zap.Error(err))
		return nil, err
	}

	var photos []string

	// Create item
	item := &models.Item{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Location:    req.Location,
		HasPhotos:   false,
		AuthorID:    req.AuthorID,
		CategoryID:  req.CategoryID,
	}

	if req.Photos != nil {
		photos = req.Photos
	}

	if err := s.itemRepo.Create(item, photos); err != nil {
		s.logger.Error("Failed to create item", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Item created successfully", zap.Int("item_id", item.ID))
	return item, nil
}

// GetItemByID retrieves an item by ID
func (s *ItemService) GetItemByID(id int) (*models.ItemResponse, error) {
	item, err := s.itemRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get item by ID", zap.Int("item_id", id), zap.Error(err))
		return nil, err
	}

	if item == nil {
		return nil, sql.ErrNoRows
	}

	return item, nil
}

// GetAllItems retrieves all items with optional filtering
func (s *ItemService) GetAllItems(filter *models.ItemFilter) ([]models.ItemResponse, error) {
	// Set default pagination if not provided
	if filter != nil {
		if filter.Limit <= 0 {
			filter.Limit = 20 // Default limit
		}
		if filter.Offset < 0 {
			filter.Offset = 0
		}
	}

	items, err := s.itemRepo.GetAll(filter)
	if err != nil {
		s.logger.Error("Failed to get all items", zap.Error(err))
		return nil, err
	}

	return items, nil
}

// UpdateItem updates an item
func (s *ItemService) UpdateItem(id int, req *models.UpdateItemRequest) (*models.ItemResponse, error) {
	if err := req.Validate(); err != nil {
		s.logger.Error("Invalid update item request", zap.Error(err))
		return nil, err
	}

	currentItemResponse, err := s.itemRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get item for update", zap.Int("item_id", id), zap.Error(err))
		return nil, err
	}
	if currentItemResponse == nil {
		return nil, sql.ErrNoRows
	}

	itemToUpdate := &models.Item{
		ID:          currentItemResponse.ID,
		Title:       currentItemResponse.Title,
		Description: currentItemResponse.Description,
		Price:       currentItemResponse.Price,
		Location:    currentItemResponse.Location,
		HasPhotos:   currentItemResponse.HasPhotos,
		AuthorID:    currentItemResponse.AuthorID,
		CategoryID:  currentItemResponse.Category.ID,
		CreatedAt:   currentItemResponse.CreatedAt,
	}

	if req.Title != nil {
		itemToUpdate.Title = *req.Title
	}
	if req.Description != nil {
		itemToUpdate.Description = *req.Description
	}
	if req.Price != nil {
		itemToUpdate.Price = *req.Price
	}
	if req.Location != nil {
		itemToUpdate.Location = *req.Location
	}
	if req.CategoryID != nil {
		itemToUpdate.CategoryID = *req.CategoryID
	}

	if err := s.itemRepo.Update(itemToUpdate); err != nil {
		s.logger.Error("Failed to update item", zap.Int("item_id", id), zap.Error(err))
		return nil, err
	}

	// 5. Снова получаем обновленные данные в формате ItemResponse, чтобы вернуть их клиенту
	updatedItem, err := s.itemRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get updated item after update", zap.Int("item_id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Item updated successfully", zap.Int("item_id", id))
	return updatedItem, nil
}

// DeleteItem deletes an item
func (s *ItemService) DeleteItem(id int) error {
	// Check if item exists
	item, err := s.itemRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get item for deletion", zap.Int("item_id", id), zap.Error(err))
		return err
	}

	if item == nil {
		return sql.ErrNoRows
	}

	// Delete item
	if err := s.itemRepo.Delete(id); err != nil {
		s.logger.Error("Failed to delete item", zap.Int("item_id", id), zap.Error(err))
		return err
	}

	s.logger.Info("Item deleted successfully", zap.Int("item_id", id))
	return nil
}

// GetItemsByLocation retrieves items by location
func (s *ItemService) GetItemsByLocation(location string) ([]models.ItemResponse, error) {
	if location == "" {
		return nil, errors.New("location cannot be empty")
	}

	items, err := s.itemRepo.GetByLocation(location)
	if err != nil {
		s.logger.Error("Failed to get items by location", zap.String("location", location), zap.Error(err))
		return nil, err
	}

	return items, nil
}

// GetAvailableItems retrieves only available items
func (s *ItemService) GetAvailableItems() ([]models.ItemResponse, error) {
	items, err := s.itemRepo.GetAvailableItems()
	if err != nil {
		s.logger.Error("Failed to get available items", zap.Error(err))
		return nil, err
	}

	return items, nil
}

// GetItemsByCategory retrieves items by category
func (s *ItemService) GetItemsByCategory(categoryID int) ([]models.ItemResponse, error) {
	if categoryID <= 0 {
		return nil, errors.New("category_id must be greater than 0")
	}

	items, err := s.itemRepo.GetByCategory(categoryID)
	if err != nil {
		s.logger.Error("Failed to get items by category", zap.Int("category_id", categoryID), zap.Error(err))
		return nil, err
	}

	return items, nil
}
