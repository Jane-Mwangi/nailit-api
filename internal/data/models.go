package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrEditConflict         = errors.New("edit conflict")
	ErrDuplicateService     = errors.New("service with this name already exists")
	ErrDuplicateServiceType = errors.New("duplicate service type")
)

type Models struct {
	Services     *ServiceModel
	ServiceTypes *ServiceTypesModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Services:     &ServiceModel{DB: db},
		ServiceTypes: &ServiceTypesModel{DB: db},
	}
}
