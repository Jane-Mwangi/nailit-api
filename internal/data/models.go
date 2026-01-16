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
	ErrDuplicateStaff       = errors.New("duplicate staff")
)

type Models struct {
	Permissions  *PermissionModel
	Services     *ServiceModel
	ServiceTypes *ServiceTypesModel
	Staff        *StaffModel
	Tokens       *TokenModel
	Users        *UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Permissions:  &PermissionModel{DB: db},
		Services:     &ServiceModel{DB: db},
		ServiceTypes: &ServiceTypesModel{DB: db},
		Staff:        &StaffModel{DB: db},
		Tokens:       &TokenModel{DB: db},
		Users:        &UserModel{DB: db},
	}
}
