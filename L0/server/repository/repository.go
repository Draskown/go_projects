package repository

import (
	"github.com/Draskown/WBL0/model"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/stan.go"
)

// Interface for displaying the order on the database level
type DBConv interface {
	ShowOrder(id string) (model.Order, error)
	ConnectStan(cfg StanCfg) (stan.Conn, stan.Subscription, error)
}

// Repository structure to hold the interface
type Repository struct {
	DBConv
}

// Creates a new repository dependant on the database itself
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		DBConv: NewDBConvPostgres(db),
	}
}
