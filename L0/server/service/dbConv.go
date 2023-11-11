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
func (s *DBConvService) ShowOrder(order model.Order) (int, error) {
	return s.repo.ShowOrder(order)
}

// DEBUG
func (s *DBConvService) TestGetDB(id int) (model.Test, error) {
	return s.repo.TestGetDB(id)
}

func (s *DBConvService) TestPostDB(test model.Test) (int, error) {
	return s.repo.TestPostDB(test)
}

func (s *DBConvService) ShowTestDB() ([]model.Test, error) {
	return s.repo.ShowTestDB()
}

func (s *DBConvService) ShowTestDBbyId(id int) (model.Test, error) {
	return s.repo.ShowTestDBbyId(id)
}
