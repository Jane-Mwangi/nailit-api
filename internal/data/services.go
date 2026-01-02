package data

import (
	"context"
	"database/sql"
	"errors"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&service.ID, &service.CreatedAt,
	)
}

func (m ServiceModel) Get(id uuid.UUID) (*Service, error) {

	query := `
 SELECT id, created_at, name
 FROM services
 WHERE id = $1`

	// declare a Service struct to hold the data returned by the query

	var service Service

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&service.ID,
		&service.CreatedAt,
		&service.Name,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &service, nil
}
