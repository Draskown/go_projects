package repository

import (
	"errors"
	"fmt"

	"github.com/Draskown/WBL0/model"
	"github.com/jmoiron/sqlx"
)

// Struct for implementing the DBConv interface
type DBConvPostgres struct {
	db *sqlx.DB
}

// Creates new struct with a provided db
func NewDBConvPostgres(db *sqlx.DB) *DBConvPostgres {
	return &DBConvPostgres{db: db}
}

// ShowOrder's core implementation
func (r *DBConvPostgres) ShowOrder(id string) (model.Order, error) {
	for _, order := range Msgs {
		if order.OrderId == id {
			return order, nil
		}
	}

	return model.Order{}, errors.New(fmt.Sprintf("No entry for id (%s)", id))
}
