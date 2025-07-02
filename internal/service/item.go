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

	// Create item
	item := &models.Item{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Location:    req.Location,
		AuthorID:    req.AuthorID,
		CategoryID:  req.CategoryID,
	}

	if err := s.itemRepo.Create(item); err != nil {
		s.logger.Error("Failed to create item", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Item created successfully", zap.Int("item_id", item.ID))
	return item, nil
}

// GetItemByID retrieves an item by ID
func (s *ItemService) GetItemByID(id int) (*models.Item, error) {
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
func (s *ItemService) GetAllItems(filter *models.ItemFilter) ([]models.Item, error) {
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
func (s *ItemService) UpdateItem(id int, req *models.UpdateItemRequest) (*models.Item, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		s.logger.Error("Invalid update item request", zap.Error(err))
		return nil, err
	}

	// Get existing item
	item, err := s.itemRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get item for update", zap.Int("item_id", id), zap.Error(err))
		return nil, err
	}

	if item == nil {
		return nil, sql.ErrNoRows
	}

	// Update fields if provided
	if req.Title != nil {
		item.Title = *req.Title
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.Price != nil {
		item.Price = *req.Price
	}
	if req.Location != nil {
		item.Location = *req.Location
	}
	if req.CategoryID != nil {
		item.CategoryID = *req.CategoryID
	}

	// Update item in database
	if err := s.itemRepo.Update(item); err != nil {
		s.logger.Error("Failed to update item", zap.Int("item_id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Item updated successfully", zap.Int("item_id", id))
	return item, nil
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
func (s *ItemService) GetItemsByLocation(location string) ([]models.Item, error) {
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
func (s *ItemService) GetAvailableItems() ([]models.Item, error) {
	items, err := s.itemRepo.GetAvailableItems()
	if err != nil {
		s.logger.Error("Failed to get available items", zap.Error(err))
		return nil, err
	}

	return items, nil
}
