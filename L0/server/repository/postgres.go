package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	deliveriesTable = "deliveries"
	ItemsTable = "items"
	paymentsTable = "payment"
	ordersTable = "orders"
)

// Struct to hold db info for simplicity of application's structure
type DBConfig struct {
	Host string
	Port string
	Username string
	Password string
	DBName string
	SSLMode string
}

// Connects to a new db
func NewPostgresDB(cfg DBConfig) (*sqlx.DB, error) {
	// Open connection with provided parameters
	db, err := sqlx.Open("postgres", 
				fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
					cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	// Ping the connection
	if pingErr := db.Ping(); pingErr != nil {
		return nil, pingErr
	}

	return db, nil
}