package service

import (
	"github.com/Draskown/WBL0/model"
	"github.com/Draskown/WBL0/server/repository"
)

// Interface for implementing DBConv service
type DBConv interface {
	ShowOrder(id string) (model.Order, error)
}

// Structure of the service to realise the interfaces
type Service struct {
	DBConv
}

// Creates services dependant on the repository
func NewService(repo *repository.Repository) *Service {
	return &Service{
		DBConv: NewDBConvService(repo.DBConv),
	}
}
