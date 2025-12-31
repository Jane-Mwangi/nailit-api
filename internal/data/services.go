package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/Jane-Mwangi/nailit-api/internal/validator"
	"github.com/google/uuid"
)

type Service struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func ValidateService(v *validator.Validator, service *Service) {
	v.Check(service.Name != "", "Name", "must be provided")

}

type ServiceModel struct {
	DB *sql.DB
}

func (m *ServiceModel) Insert(service *Service) error {

	query := `
        INSERT INTO services (name)
        VALUES ($1)
        RETURNING id, created_at
    `
	args := []interface{}{
		service.Name,
	}

	// Create a context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the QueryRow() method to execute the query and scan the result into the movie struct.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&service.ID, &service.CreatedAt,
	)
}
