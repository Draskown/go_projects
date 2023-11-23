package service

import (
	"github.com/Draskown/WBL0/model"
	"github.com/Draskown/WBL0/server/repository"
)

// Structure to hold a
// service interface dependant on the repository
// on the service level
type DBConvService struct {
	repo repository.DBConv
}

// Creates a new DBConv service from the repo's interface
func NewDBConvService(repo repository.DBConv) *DBConvService {
	return &DBConvService{repo: repo}
}

// Implements the service interface for calling from the service level
func (s *DBConvService) ShowOrder(id string) (model.Order, error) {
	return s.repo.ShowOrder(id)
}
