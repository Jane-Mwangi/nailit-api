package data

import "database/sql"

type Models struct {
	Services *ServiceModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Services: &ServiceModel{DB: db},
	}
}
