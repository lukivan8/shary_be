package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"shary_be/internal/models"

	"github.com/jmoiron/sqlx"
)

// RentRepository handles database operations for rents
type RentRepository struct {
	db *sqlx.DB
}

// NewRentRepository creates a new rent repository
func NewRentRepository(db *sqlx.DB) *RentRepository {
	return &RentRepository{db: db}
}

// GetAll retrieves all rents with optional filtering
func (r *RentRepository) GetAll(filter *models.RentFilter) ([]models.Rent, error) {
	var rents []models.Rent

	var queryBuilder strings.Builder

	// Build dynamic query with filters
	queryBuilder.WriteString(`
        SELECT
            s.id, s.item_id, s.user_id, s.start_date, s.end_date, s.price, s.status_id, s.created_at, s.updated_at,
            c.id AS "status.id", c.name AS "status.name"
        FROM rents s
        LEFT JOIN statuses c ON s.status_id = c.id
        WHERE 1=1
    `)

	var args []interface{}

	// Add filters
	if filter != nil {
		if filter.MinStartDate != nil {
			queryBuilder.WriteString(" AND r.start_date >= ?")
			args = append(args, *filter.MinStartDate)
		}
		if filter.MaxStartDate != nil {
			queryBuilder.WriteString(" AND r.start_date <= ?")
			args = append(args, *filter.MaxStartDate)
		}
		if filter.MinEndDate != nil {
			queryBuilder.WriteString(" AND r.end_date >= ?")
			args = append(args, *filter.MinEndDate)
		}
		if filter.MaxEndDate != nil {
			queryBuilder.WriteString(" AND r.end_date <= ?")
			args = append(args, *filter.MaxEndDate)
		}
		if filter.Price != nil {
			queryBuilder.WriteString(" AND r.price >= ?")
			args = append(args, *filter.Price)
		}
		if filter.StatusID != nil {
			queryBuilder.WriteString(" AND r.status_id = ?")
			args = append(args, *filter.StatusID)
		}
	}

	// Add ordering
	queryBuilder.WriteString(" ORDER BY r.created_at DESC")

	// Add pagination
	if filter != nil {
		if filter.Limit > 0 {
			queryBuilder.WriteString(" LIMIT ?")
			args = append(args, filter.Limit)
		}
		if filter.Offset > 0 {
			queryBuilder.WriteString(" OFFSET ?")
			args = append(args, filter.Offset)
		}
	}

	query := r.db.Rebind(queryBuilder.String())

	err := r.db.Select(&rents, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get all rents with filter: %w", err)
	}

	return rents, nil
}

// Create creates a new rent in the database
func (r *RentRepository) Create(rent *models.Rent) error {
	query := `
		INSERT INTO rents (item_id, user_id, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	now := time.Now()
	rent.CreatedAt = now
	rent.UpdatedAt = now

	return r.db.QueryRow(
		query,
		rent.ItemID,
		rent.UserID,
		rent.StartDate,
		rent.EndDate,
		rent.CreatedAt,
		rent.UpdatedAt,
	).Scan(&rent.ID)
}

// GetByItemID retrieves rents by item ID
func (r *RentRepository) GetByItemID(itemID int) ([]models.Rent, error) {
	var rents []models.Rent
	query := `SELECT * FROM rents WHERE item_id = $1`

	err := r.db.Select(&rents, query, itemID)
	if err != nil {
		return nil, err
	}

	return rents, nil
}

// GetByUserID retrieves rents by user ID
func (r *RentRepository) GetByUserID(userID int) ([]models.Rent, error) {
	var rents []models.Rent
	query := `SELECT * FROM rents WHERE user_id = $1`

	err := r.db.Select(&rents, query, userID)
	if err != nil {
		return nil, err
	}

	return rents, nil
}

// GetByID retrieves a rent by ID
func (r *RentRepository) GetByID(id int) (*models.Rent, error) {
	var rent models.Rent
	query := `SELECT * FROM rents WHERE id = $1`

	err := r.db.Get(&rent, query, id)
	if err != nil {
		return nil, err
	}

	return &rent, nil
}

// Update updates an existing rent in the database
func (r *RentRepository) Update(tx *sqlx.Tx, rent *models.RentToUpdate) error {
	query := `
		UPDATE rents
		SET start_date = $1, end_date = $2, updated_at = $3
		WHERE id = $4`

	now := time.Now()
	rent.UpdatedAt = now

	result, err := tx.Exec(query,
		rent.StartDate,
		rent.EndDate,
		rent.UpdatedAt,
		rent.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete deletes a rent by ID
func (r *RentRepository) Delete(id int) error {
	query := `DELETE FROM rents WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetAvailableRents retrieves only available rents
func (r *RentRepository) GetAvailableRents() ([]models.Rent, error) {
	var rents []models.Rent
	query := `SELECT * FROM rents WHERE end_date > NOW()`

	err := r.db.Select(&rents, query)
	if err != nil {
		return nil, err
	}

	return rents, nil
}


