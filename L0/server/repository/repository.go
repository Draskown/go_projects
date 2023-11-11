package repository

import (
	"github.com/Draskown/WBL0/model"
	"github.com/jmoiron/sqlx"
)

// Interface for displaying the order on the database level
type DBConv interface {
	ShowOrder(order model.Order) (int, error)

	// DEBUG
	TestGetDB(id int) (model.Test, error)
	TestPostDB(test model.Test) (int, error)
	ShowTestDB() ([]model.Test, error)
	ShowTestDBbyId(id int) (model.Test, error)
}

// Repository structure to hold the interface
type Repository struct {
	DBConv
}

// Creates a new repository dependant on the database itself
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		DBConv: NewDBConvRepo(db),
	}
}
