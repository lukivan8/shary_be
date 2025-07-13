package service

import (
	"shary_be/internal/models"
	"shary_be/internal/repository"

	"go.uber.org/zap"
)

type ItemPhotoService struct {
	itemPhotoRepo *repository.ItemPhotoRepository
	logger        *zap.Logger
}

func NewItemPhotoService(itemPhotoRepo *repository.ItemPhotoRepository, logger *zap.Logger) *ItemPhotoService {
	return &ItemPhotoService{
		itemPhotoRepo: itemPhotoRepo,
		logger:        logger,
	}
}

// GetPhotosByItemID retrieves all photos for an item
func (s *ItemPhotoService) GetPhotosByItemID(itemID int) ([]models.ItemPhoto, error) {
	photos, err := s.itemPhotoRepo.GetPhotosByItemID(itemID)
	if err != nil {
		s.logger.Error("Failed to get photos by item ID", zap.Int("item_id", itemID), zap.Error(err))
		return nil, err
	}

	return photos, nil
}

func (s *ItemPhotoService) AddPhotos(itemID int, photoURLs []string) error {
	if len(photoURLs) == 0 {
		return nil
	}

	for _, url := range photoURLs {
		if err := s.itemPhotoRepo.Add(itemID, url); err != nil {
			s.logger.Error("Failed to add photo", zap.Int("item_id", itemID), zap.String("url", url), zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *ItemPhotoService) DeletePhotos(itemID int, photoIDs []int) error {
	if len(photoIDs) == 0 {
		return nil
	}

	for _, id := range photoIDs {
		if err := s.itemPhotoRepo.Delete(id); err != nil {
			s.logger.Error("Failed to delete photo", zap.Int("item_id", itemID), zap.Int("photo_id", id), zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *ItemPhotoService) CountPhotosByItemID(itemID int) (int, error) {
	count, err := s.itemPhotoRepo.CountByItemID(itemID)
	if err != nil {
		s.logger.Error("Failed to count photos", zap.Int("item_id", itemID), zap.Error(err))
		return 0, err
	}

	return count, nil
}
