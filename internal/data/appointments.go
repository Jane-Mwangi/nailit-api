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

type AppointmentStatus string

const (
	StatusPending   AppointmentStatus = "pending"
	StatusBooked    AppointmentStatus = "booked"
	StatusCancelled AppointmentStatus = "cancelled"
	StatusCompleted AppointmentStatus = "completed"
)

type Appointment struct {
	ID            uuid.UUID         `json:"id"`
	CustomerID    uuid.UUID         `json:"customer_id"`
	ServiceID     uuid.UUID         `json:"service_id"`
	ServiceTypeID uuid.UUID         `json:"service_type_id"`
	StartsAt      time.Time         `json:"starts_at"`
	EndsAt        time.Time         `json:"ends_at"`
	Status        AppointmentStatus `json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

func ValidateAppointment(v *validator.Validator, a *Appointment) {

	v.Check(a.ServiceTypeID != uuid.Nil, "service_type_id", "must be provided")
	v.Check(!a.StartsAt.IsZero(), "starts_at", "must be provided")
	v.Check(a.EndsAt.After(a.StartsAt), "ends_at", "must be be after starts_at")

	duration := a.EndsAt.Sub(a.StartsAt)
	v.Check(duration > 0, "duration", "must be positive")
	v.Check(duration <= 8*time.Hour, "duration", "must not exceed 8 hours")
}

type AppointmentModel struct {
	DB *sql.DB
}

func (a *AppointmentModel) Insert(appointment *Appointment) error {

	query := `
        INSERT INTO appointments (
	        customer_id,
			service_type_id,
			starts_at,
			ends_at,
			status
		)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `
	args := []interface{}{
		appointment.CustomerID,
		appointment.ServiceTypeID,
		appointment.StartsAt,
		appointment.EndsAt,
		appointment.Status,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, args...).Scan(
		&appointment.ID, &appointment.CreatedAt, &appointment.UpdatedAt,
	)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {

			switch pgErr.Code {

			case "23P01":
				return ErrOverlappingAppointment

			case "23505":
				return ErrDuplicateAppointment
			}
		}

		return err
	}

	return nil
}

func (a *AppointmentModel) Get(id uuid.UUID) (*Appointment, error) {

	query := `
		SELECT
    id,
    customer_id,
    service_type_id,
    starts_at,
    ends_at,
    status,
    created_at,
    updated_at
FROM appointments
WHERE id = $1
	`

	var appointment Appointment

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, id).Scan(
		&appointment.ID,
		&appointment.CustomerID,
		&appointment.ServiceTypeID,
		&appointment.StartsAt,
		&appointment.EndsAt,
		&appointment.Status,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &appointment, nil
}

func (a *AppointmentModel) Update(appointment *Appointment) error {

	query := `
    UPDATE appointments
    SET
	   service_type_id = $1,
       starts_at = $2,
       ends_at = $3,
       status = $4,
       updated_at = NOW()
    WHERE id = $5
    RETURNING updated_at`

	args := []any{
		appointment.ServiceTypeID,
		appointment.StartsAt,
		appointment.EndsAt,
		appointment.Status,
		appointment.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, args...).Scan(
		&appointment.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (a *AppointmentModel) Delete(id uuid.UUID) error {

	query := `
		DELETE FROM appointments
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := a.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (a *AppointmentModel) GetAll(filters Filters) ([]*Appointment, Metadata, error) {

	query := fmt.Sprintf(`
		SELECT count(*) OVER(),
		       id, customer_id, service_type_id,
		       starts_at, ends_at, status, created_at, updated_at
		FROM appointments
		ORDER BY %s %s, id ASC
		LIMIT $1 OFFSET $2
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		filters.limit(),
		filters.offset(),
	}

	rows, err := a.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var totalRecords int
	var appointments []*Appointment

	for rows.Next() {

		var appt Appointment

		err := rows.Scan(
			&totalRecords,
			&appt.ID,
			&appt.CustomerID,
			&appt.ServiceTypeID,
			&appt.StartsAt,
			&appt.EndsAt,
			&appt.Status,
			&appt.CreatedAt,
			&appt.UpdatedAt,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		appointments = append(appointments, &appt)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return appointments, metadata, nil
}
