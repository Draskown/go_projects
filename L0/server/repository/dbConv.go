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

func (r *DBConvPostgres) insertOrder(order model.Order) error {
	// Initialise DB transaction
	tx, err := r.db.Begin()
	if err != nil {
		return errors.New(fmt.Sprintf("Could not begin transaction (%s)\n", err.Error()))
	}

	// Insert message's values into the
	// delivery table
	query := fmt.Sprintf(`
	INSERT INTO %s
	(name, phone, zip, city, address, region, email)
	VALUES
	($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`, deliveriesTable)

	// Get the delivery id and assign it
	// to the order struct's delivery object
	row := tx.QueryRow(query,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.ZipCode,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err := row.Scan(&order.Delivery.Id); err != nil {
		tx.Rollback()
		return errors.New(fmt.Sprintf("Could not insert into deliveries table (%s)\n", err.Error()))
	}

	// Range over all items that order has
	for _, item := range order.Items {
		// Insert item into the items table
		query = fmt.Sprintf(`
		INSERT INTO %s
		(track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, chrt_id)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`, itemsTable)

		// Get item id
		row = tx.QueryRow(query,
			item.TrackNumber,
			item.Price,
			item.RId,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmId,
			item.Brand,
			item.Status,
			item.ChartId,
		)
		if err := row.Scan(&item.Id); err != nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("Could not insert into items table (%s)\n", err.Error()))
		}

		// Add it to the itemIds array of ints
		order.ItemsIds = append(order.ItemsIds, int64(item.Id))
	}

	// Insert order into the payments table
	query = fmt.Sprintf(`
	INSERT INTO %s
	(transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`, paymentsTable)

	// Get the payment id and assign it
	// to the order struct's payment object
	row = tx.QueryRow(query,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err := row.Scan(&order.Payment.Id); err != nil {
		tx.Rollback()
		return errors.New(fmt.Sprintf("Could not insert into payments table (%s)\n", err.Error()))
	}

	// Insert the order itself
	query = fmt.Sprintf(`
	INSERT INTO %s 
	(order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING id`, ordersTable)

	row = tx.QueryRow(query,
		order.OrderId,
		order.TrackNumber,
		order.Entry,
		order.Delivery.Id,
		order.Payment.Id,
		order.ItemsIds,
		order.Locale,
		order.InternalSignature,
		order.CustomerId,
		order.DeliveryService,
		order.ShardKey,
		order.SmId,
		order.DateCreated,
		order.OofShard,
	)

	if err := row.Scan(&order.Id); err != nil {
		tx.Rollback()
		return errors.New(fmt.Sprintf("Could not insert into orders table (%s)\n", err.Error()))
	}

	// Finish the transaction
	return tx.Commit()
}

// Restores Msgs cache from DB
func (r *DBConvPostgres) restoreCache() error {
	var orders []model.Order

	// Run a joined query to write two structs of the order
	// (excepting items as they are handled differently)
	query := fmt.Sprintf(`
	SELECT o.id, o.order_uid, o.track_number, o.entry, o.items, o.locale, 
		o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, 
		o.date_created, o.oof_shard, 
		d.id "delivery.id", d.name "delivery.name", d.phone "delivery.phone", 
		d.zip "delivery.zip", d.city "delivery.city", d.address "delivery.address", 
		d.region "delivery.region", d.email "delivery.email", 
		p.id "payment.id", p.transaction "payment.transaction", p.request_id "payment.request_id", 
		p.currency "payment.currency", p.provider "payment.provider", p.amount "payment.amount", 
		p.payment_dt "payment.payment_dt", p.bank "payment.bank", p.delivery_cost "payment.delivery_cost", 
		p.goods_total "payment.goods_total", p.custom_fee "payment.custom_fee"
	FROM %s o
	JOIN %s d ON o.delivery=d.id
	JOIN %s p ON o.payment=p.id
	`, ordersTable, deliveriesTable, paymentsTable)

	if err := r.db.Select(&orders, query); err != nil {
		return err
	}

	// Range over items in all the orders
	for _, order := range orders {
		for _, id := range order.ItemsIds {
			var item model.Item

			// Get the item from DB by its id in the order slice
			query = fmt.Sprintf(`
			SELECT * FROM %s
			WHERE id=$1
			`, itemsTable)

			if err := r.db.Get(&item, query, id); err != nil {
				return err
			}

			order.Items = append(order.Items, item)
		}
		// Append the filled out order to the cache
		Msgs = append(Msgs, order)
	}

	return nil
}
