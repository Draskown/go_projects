package repository

import (
	"errors"
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
func (r *DBConvPostgres) ShowOrder(order model.Order) (int, error) {
	var id int

	query := fmt.Sprintf(`
	INSERT INTO %s 
	(order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING id`, ordersTable)

	row := r.db.QueryRow(query,
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
func (r *DBConvPostgres) TestGetDB(id int) (model.Test, error) {
	var test model.Test

	query := fmt.Sprintf(`
	SELECT * FROM %s
	WHERE id = $1`, "tests")

	err := r.db.Get(&test, query, id)
	if err != nil {
		return model.Test{}, err
	}

	return test, nil
}

func (r *DBConvPostgres) TestPostDB(test model.Test) (int, error) {
	var id int

	query := fmt.Sprintf(`
	INSERT INTO %s
	(value, text, arr) VALUES
	($1, $2, $3)
	RETURNING id`, "tests")

	row := r.db.QueryRow(query,
		test.Value, test.Text, test.Arr)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *DBConvPostgres) ShowTestDB() ([]model.Test, error) {
	var result []model.Test
	query := fmt.Sprintf(`SELECT * FROM %s;`, "tests")

	if err := r.db.Select(&result, query); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("Empty return from DB!")
	}

	return result, nil
}

func (r *DBConvPostgres) ShowTestDBbyId(id int) (model.Test, error) {
	var result model.Test

	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1;`, "tests")

	if err := r.db.Get(&result, query, id); err != nil {
		return model.Test{}, err
	}

	return result, nil
}
