package service

import (
	"github.com/Draskown/WBL0/model"
	"github.com/Draskown/WBL0/server/repository"
)

// Interface for implementing the service
type DBConv interface {
	ShowOrder(order model.Order) (int, error)

	// DEBUG
	TestGetDB (id int) (model.Test, error)
	TestPostDB (test model.Test) (int, error)
}

// Structure of the service that contains the interface
type Service struct{
	DBConv
}

// Creates a new service from implementing
// repository's (meaning it is handled on the database level)
// interface
func NewService(repo *repository.Repository) *Service {
	return &Service{
		DBConv: NewDBConvService(repo.DBConv),
	}
}