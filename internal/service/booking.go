package service

import (
	"database/sql"

	"shary_be/internal/models"
	"shary_be/internal/repository"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type RentService struct {
	rentRepo *repository.RentRepository
	logger      *zap.Logger
	db          *sqlx.DB
}

func NewRentService(rentRepo *repository.RentRepository, logger *zap.Logger, db *sqlx.DB) *RentService {
	return &RentService{
		rentRepo: rentRepo,
		logger:      logger,
		db:          db,
	}
}

// GetAllRents retrieves all rents with optional filtering
func (s *RentService) GetAllRents(filter *models.RentFilter) ([]models.Rent, error) {
	// Set default pagination if not provided
	if filter != nil {
		if filter.Limit <= 0 {
			filter.Limit = 20 // Default limit
		}
		if filter.Offset < 0 {
			filter.Offset = 0
		}
	}

	rents, err := s.rentRepo.GetAll(filter)
	if err != nil {
		s.logger.Error("Failed to get all rents", zap.Error(err))
		return nil, err
	}

	return rents, nil
}

// CreateRent creates a new rent for the given item and user
func (s *RentService) CreateRent(req *models.CreateRentRequest) (*models.Rent, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		s.logger.Error("Invalid create rent request", zap.Error(err))
		return nil, err
	}

	// Create rent
	rent := &models.Rent{
		ItemID: req.ItemID,
		UserID: req.UserID,
		StartDate: req.StartDate,
		EndDate: req.EndDate,
	}

	if err := s.rentRepo.Create(rent); err != nil {
		s.logger.Error("Failed to create rent", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Rent created successfully", zap.Int("rent_id", rent.ID))
	return rent, nil
}

// GetByItemID retrieves rents by item ID
func (s *RentService) GetByItemID(itemID int) ([]models.Rent, error) {
	rents, err := s.rentRepo.GetByItemID(itemID)
	if err != nil {
		s.logger.Error("Failed to get rents by item ID", zap.Int("item_id", itemID), zap.Error(err))
		return nil, err
	}

	return rents, nil
}

// GetByUserID retrieves rents by user ID
func (s *RentService) GetByUserID(userID int) ([]models.Rent, error) {
	rents, err := s.rentRepo.GetByUserID(userID)
	if err != nil {
		s.logger.Error("Failed to get rents by user ID", zap.Int("user_id", userID), zap.Error(err))
		return nil, err
	}

	return rents, nil
}

// GetByID retrieves a rent by ID
func (s *RentService) GetByID(id int) (*models.Rent, error) {
	rent, err := s.rentRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get rent by ID", zap.Int("rent_id", id), zap.Error(err))
		return nil, err
	}

	if rent == nil {
		return nil, sql.ErrNoRows
	}

	return rent, nil
}

// Update updates an existing rent
func (s *RentService) Update(id int, req *models.UpdateRentRequest) (*models.Rent, error) {
	if err := req.Validate(); err != nil {
		s.logger.Error("Invalid update rent request", zap.Error(err))
		return nil, err
	}

	// 1. Start Transaction
	tx, err := s.db.Beginx()
	if err != nil {
		s.logger.Error("Failed to begin transaction", zap.Error(err))
		return nil, err
	}

	defer tx.Rollback()

	// Get current rent data
	currentRent, err := s.rentRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get rent for update", zap.Int("rent_id", id), zap.Error(err))
		return nil, err
	}

	// 2. Update the main rent record
	rentToUpdate := &models.RentToUpdate{
		ID:          id,
		ItemID:      currentRent.ItemID,
		UserID:      currentRent.UserID,
		StartDate:   currentRent.StartDate,
		EndDate:     currentRent.EndDate,
		UpdatedAt:   currentRent.UpdatedAt,
		CreatedAt:   currentRent.CreatedAt,
	}

	if req.StartDate != nil {
		rentToUpdate.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		rentToUpdate.EndDate = *req.EndDate
	}

	if err := s.rentRepo.Update(tx, rentToUpdate); err != nil {
		s.logger.Error("Failed to update rent", zap.Int("rent_id", id), zap.Error(err))
		return nil, err
	}

	// 3. Commit the transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", zap.Int("rent_id", id), zap.Error(err))
		return nil, err
	}

	updatedRent, err := s.rentRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get updated rent after commit", zap.Int("rent_id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Rent updated successfully", zap.Int("rent_id", id))
	return updatedRent, nil
}

// Delete deletes a rent by ID
func (s *RentService) Delete(id int) error {
	// Check if rent exists
	rent, err := s.rentRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get rent for deletion", zap.Int("rent_id", id), zap.Error(err))
		return err
	}

	if rent == nil {
		return sql.ErrNoRows
	}

	// Delete rent
	if err := s.rentRepo.Delete(id); err != nil {
		s.logger.Error("Failed to delete rent", zap.Int("rent_id", id), zap.Error(err))
		return err
	}

	s.logger.Info("Rent deleted successfully", zap.Int("rent_id", id))
	return nil
}

// GetAvailableRents retrieves only available rents
func (s *RentService) GetAvailableRents() ([]models.Rent, error) {
	rents, err := s.rentRepo.GetAvailableRents()
	if err != nil {
		s.logger.Error("Failed to get available rents", zap.Error(err))
		return nil, err
	}

	return rents, nil
}