package service

import (
	"shary_be/internal/models"
	"shary_be/internal/repository"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type ItemPhotoService struct {
	itemPhotoRepo *repository.ItemPhotoRepository
	itemRepo      *repository.ItemRepository
	logger        *zap.Logger
	db            *sqlx.DB
}

func NewItemPhotoService(itemPhotoRepo *repository.ItemPhotoRepository, itemRepo *repository.ItemRepository, logger *zap.Logger, db *sqlx.DB) *ItemPhotoService {
	return &ItemPhotoService{
		itemPhotoRepo: itemPhotoRepo,
		itemRepo:      itemRepo,
		logger:        logger,
		db:            db,
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
	tx, err := s.db.Beginx()
	if err != nil {
		s.logger.Error("Failed to begin transaction", zap.Error(err))
		return err
	}

	defer tx.Rollback()

	if len(photoURLs) == 0 {
		return nil
	}

	if err := s.itemPhotoRepo.Add(tx, itemID, photoURLs); err != nil {
		s.logger.Error("Failed to add photo", zap.Int("item_id", itemID), zap.Error(err))
		return err
	}

	photoCount, err := s.itemPhotoRepo.CountByItemID(tx, itemID)
	if err != nil {
		s.logger.Error("Failed to count photos", zap.Int("item_id", itemID), zap.Error(err))
		return err
	}

	if err := s.itemRepo.UpdateHasPhotos(tx, itemID, photoCount > 0); err != nil {
		s.logger.Error("Failed to update item", zap.Int("item_id", itemID), zap.Error(err))
		return err
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", zap.Int("item_id", itemID), zap.Error(err))
		return err
	}

	s.logger.Info("Successfully added photos", zap.Int("item_id", itemID), zap.Int("photo_count", photoCount))

	return nil
}

func (s *ItemPhotoService) DeletePhotos(itemID int, photoIDs []int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		s.logger.Error("Failed to begin transaction", zap.Error(err))
		return err
	}

	defer tx.Rollback()

	if len(photoIDs) == 0 {
		return nil
	}

	if err := s.itemPhotoRepo.Delete(tx, photoIDs); err != nil {
		s.logger.Error("Failed to bulk delete photos", zap.Int("item_id", itemID), zap.Error(err))
		return err
	}

	photoCount, err := s.itemPhotoRepo.CountByItemID(tx, itemID)
	if err != nil {
		s.logger.Error("Failed to count photos", zap.Int("item_id", itemID), zap.Error(err))
		return err
	}

	if err := s.itemRepo.UpdateHasPhotos(tx, itemID, photoCount > 0); err != nil {
		s.logger.Error("Failed to update item", zap.Int("item_id", itemID), zap.Error(err))
		return err
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", zap.Int("item_id", itemID), zap.Error(err))
		return err
	}

	s.logger.Info("Successfully deleted photos", zap.Int("item_id", itemID), zap.Int("photo_count", photoCount))

	return nil
}

func (s *ItemPhotoService) CountPhotosByItemID(itemID int) (int, error) {
	count, err := s.itemPhotoRepo.CountByItemID(s.db, itemID)
	if err != nil {
		s.logger.Error("Failed to count photos", zap.Int("item_id", itemID), zap.Error(err))
		return 0, err
	}

	return count, nil
}
