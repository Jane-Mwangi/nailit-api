package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Jane-Mwangi/nailit-api/internal/validator"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Staff struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"version"`
}

// validator fields must match json
func ValidateStaff(v *validator.Validator, staff *Staff) {
	v.Check(staff.Name != "", "name", "must be provided")
	v.Check(staff.Email != "", "email", "must be provided")
	v.Check(len(staff.Name) <= 100, "name", "must not exceed 100 characters")

}

type StaffModel struct {
	DB *sql.DB
}

func (m *StaffModel) Insert(Staff *Staff) error {

	query := `
        INSERT INTO staff (
			name,
			email,
			is_active
		)
        VALUES ($1, $2, $3)
        RETURNING id, created_at,version
    `
	args := []interface{}{
		Staff.Name,
		Staff.Email,
		Staff.IsActive,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&Staff.ID, &Staff.CreatedAt, &Staff.Version,
	)

	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				return ErrDuplicateStaff
			}
		}
		return err
	}
	return nil
}

func (s StaffModel) Get(id uuid.UUID) (*Staff, error) {

	query := `
 SELECT id, name, email, is_active, created_at, version
 FROM staff
 WHERE id = $1`

	var staff Staff

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&staff.ID,
		&staff.Name,
		&staff.Email,
		&staff.IsActive,
		&staff.CreatedAt,
		&staff.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &staff, nil
}

func (s StaffModel) Update(staff *Staff) error {
	query := `
        UPDATE staff
        SET
		   name = $1,
		   email = $2,
		   is_active = $3,
           version = version + 1
        WHERE id = $4 AND version = $5
       RETURNING version`

	args := []interface{}{
		staff.Name,
		staff.Email,
		staff.IsActive,
		staff.ID,
		staff.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, args...).Scan(&staff.Version)
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

func (s StaffModel) Delete(id uuid.UUID) error {

	query := `
        DELETE FROM staff
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

func (s StaffModel) GetAll(name string, filters Filters) ([]*Staff, Metadata, error) {

	query := fmt.Sprintf(`
   SELECT count(*) OVER(), id, name, email, is_active, created_at, version
        FROM staff
		WHERE (
            to_tsvector('simple', name)
            @@ plainto_tsquery('simple', $1)
            OR $1 = ''
			
			)

        ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, filters.limit(), filters.offset()}

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	staffList := []*Staff{}
	// Use rows.Next to iterate through the rows in the resultset
	for rows.Next() {

		var staff Staff

		err := rows.Scan(
			&totalRecords,
			&staff.ID,
			&staff.Name,
			&staff.Email,
			&staff.IsActive,
			&staff.CreatedAt,
			&staff.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		staffList = append(staffList, &staff)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return staffList, metadata, nil
}
