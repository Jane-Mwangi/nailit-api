package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/Jane-Mwangi/nailit-api/internal/validator"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ServiceType struct {
	ID              uuid.UUID `json:"id"`
	ServiceID       uuid.UUID `json:"service_id"`
	Name            string    `json:"name"`
	Price           int       `json:"price"`
	DurationMinutes int       `json:"duration_minutes"`
	ImageURL        string    `json:"image_url"`
	CreatedAt       time.Time `json:"created_at"`
	Version         int       `json:"version"`
}

// validator fields must match json
func ValidateServiceType(v *validator.Validator, serviceType *ServiceType) {
	v.Check(serviceType.Name != "", "name", "must be provided")
	v.Check(len(serviceType.Name) <= 100, "name", "must not exceed 100 characters")
	v.Check(serviceType.Price > 0, "price", "must be greater than zero")
	v.Check(serviceType.DurationMinutes > 0, "duration_minutes", "must be greater than zero")
	// v.Check(serviceType.ImageURL != "", "image_url", "must be provided")
	v.Check(serviceType.Price <= 100000, "price", "must not be excessive")
	v.Check(serviceType.ServiceID != uuid.Nil, "service_id", "must be provided")

	// v.Check(strings.HasPrefix(serviceType.ImageURL, "http"), "image_url", "must be a valid URL")

}

type ServiceTypesModel struct {
	DB *sql.DB
}

func (m *ServiceTypesModel) Insert(serviceType *ServiceType) error {

	query := `
        INSERT INTO service_types (
		    service_id,
			name,
			price,
			duration_minutes,
			image_url
		)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at,version
    `
	args := []interface{}{
		serviceType.ServiceID,
		serviceType.Name,
		serviceType.Price,
		serviceType.DurationMinutes,
		serviceType.ImageURL,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&serviceType.ID, &serviceType.CreatedAt, &serviceType.Version,
	)

	if err != nil {
		// PostgreSQL unique constraint violation
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				return ErrDuplicateServiceType
			}
		}
		return err
	}
	return nil
}
