package data

import (
	"context"
	"database/sql"
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
