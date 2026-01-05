package data

import (
	"context"
	"database/sql"
	"errors"
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
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				return ErrDuplicateServiceType
			}
		}
		return err
	}
	return nil
}

func (s ServiceTypesModel) Get(id uuid.UUID) (*ServiceType, error) {

	query := `
 SELECT id, service_id, name, price, duration_minutes, image_url, created_at, version
 FROM service_types
 WHERE id = $1`

	// declare a Service Type struct to hold the data returned by the query

	var service_type ServiceType

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&service_type.ID,
		&service_type.ServiceID,
		&service_type.Name,
		&service_type.Price,
		&service_type.DurationMinutes,
		&service_type.ImageURL,
		&service_type.CreatedAt,
		&service_type.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &service_type, nil
}

func (s ServiceTypesModel) Update(serviceType *ServiceType) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
        UPDATE service_types
        SET
		   name = $1,
		   price = $2,
		   duration_minutes = $3,
		   image_url = $4,
           version = version + 1
        WHERE id = $5 AND version = $6
       RETURNING version`

	args := []interface{}{
		serviceType.Name,
		serviceType.Price,
		serviceType.DurationMinutes,
		serviceType.ImageURL,
		serviceType.ID,
		serviceType.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, args...).Scan(&serviceType.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (s ServiceTypesModel) Delete(id uuid.UUID) error {

	query := `
        DELETE FROM service_types
        WHERE id = $1`

	result, err := s.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
