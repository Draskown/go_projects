package repository

import (
	"fmt"

	"github.com/Draskown/WBL0/model"
	"github.com/jmoiron/sqlx"
)

// Struct to link the ShowOrder method to
type DBConvPostgres struct {
	db *sqlx.DB
}

// Create new struct with a provided db
func NewDBConvRepo(db *sqlx.DB) *DBConvPostgres {
	return &DBConvPostgres{db: db}
}

// ShowOrder's core implementation
func (s *DBConvPostgres) ShowOrder(order model.Order) (int, error) {
	var id int
	
	query := fmt.Sprintf(`
	INSERT INTO %s 
	(order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING id`, ordersTable)

	row := s.db.QueryRow(query,
		order.CustomerId,
		order.DateCreated,
		order.Delivery,
		order.DeliveryService,
		order.Entry,
		order.InternalSignature,
		order.Items,
		order.Locale,
		order.OofShard,
		order.OrderId,
		order.Payment,
		order.ShardKey,
		order.SmId,
		order.TrackNumber,
	)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	
	return id, nil
}

// DEBUG
func (s *DBConvPostgres) TestGetDB (id int) (model.Test, error) {
	// var test model.Test

	// query := fmt.Sprintf(`
	// SELECT value, text, arr
	// FROM %s WHERE id = %d`, "tests", id)
	var test model.Test
	
	query := fmt.Sprintf(`
	SELECT value, text, arr_one
	FROM %s WHERE id = %d`, "tests", id)

	rows, err := s.db.Queryx(query)
	if err != nil {
		return model.Test{}, err
	}
	
	for rows.Next() {
		if err := rows.StructScan(&test); err != nil {
			return model.Test{}, err
		}
	}
	
	return test, nil
}

func (s *DBConvPostgres) TestPostDB (test model.Test) (int, error) {
	var id int

	query := fmt.Sprintf(`
	INSERT INTO %s
	(value, text, arr_one) VALUES
	($1, $2, $3)
	RETURNING id`, "tests")

	row := s.db.QueryRow(query,
		test.Value, test.Text, test.Arr)
	
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	
	return id, nil
}