package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Jane-Mwangi/nailit-api/internal/validator"
	"github.com/google/uuid"
)

type Service struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"version"`
}

func ValidateService(v *validator.Validator, service *Service) {
	v.Check(service.Name != "", "Name", "must be provided")

}

type ServiceModel struct {
	DB *sql.DB
}

func (s *ServiceModel) Insert(service *Service) error {

	query := `
        INSERT INTO services (name)
        VALUES ($1)
        RETURNING id, created_at,version
    `
	args := []interface{}{
		service.Name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, args...).Scan(
		&service.ID, &service.CreatedAt, &service.Version,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return ErrDuplicateService
		}
		return err
	}

	return nil
}

func (s ServiceModel) Get(id uuid.UUID) (*Service, error) {

	query := `
 SELECT id, created_at, name,version
 FROM services
 WHERE id = $1`

	// declare a Service struct to hold the data returned by the query

	var service Service

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&service.ID,
		&service.CreatedAt,
		&service.Name,
		&service.Version,
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

func (s ServiceModel) Update(service *Service) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
        UPDATE services
        SET name = $1,
		version = version + 1
        WHERE id = $2 AND version = $3
       RETURNING id, name, created_at, version`

	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		service.Name,
		service.ID,
		service.Version,
	}

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the SQL query. If no matching row could be found, we know the service
	// version has changed (or the record has been deleted) and we return our custom
	// ErrEditConflict error.
	err := s.DB.QueryRowContext(ctx, query, args...).Scan(
		&service.ID,
		&service.Name,
		&service.CreatedAt,
		&service.Version,
	)
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

func (s ServiceModel) Delete(id uuid.UUID) error {

	// Construct the SQL query to delete the record.
	query := `
        DELETE FROM services
        WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := s.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the services table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m ServiceModel) GetAll(name string, filters Filters) ([]*Service, error) {
	// Construct the SQL query to retrieve all services
	query := `
        SELECT id, created_at, name,version
        FROM services
        ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	services := []*Service{}
	// Use rows.Next to iterate through the rows in the resultset
	for rows.Next() {

		var service Service

		err := rows.Scan(
			&service.ID,
			&service.CreatedAt,
			&service.Name,
			&service.Version,
		)
		if err != nil {
			return nil, err
		}

		services = append(services, &service)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}
