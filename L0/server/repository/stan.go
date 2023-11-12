package repository

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Draskown/WBL0/model"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

// Cache slice
var Msgs []model.Order

// Stan Connection configuration struct
// containing client_id`, `cluster_id`
// and `subject` to subscribe to
type StanConn struct {
	ClientId  string
	ClusterId string
	Subject   string
	DB        *sqlx.DB
}

// Creates a new connection with Stan using info provided in the StanConn struct
// , returning connection interface, subscription interface and error
func NewStanConn(cfg StanConn) (stan.Conn, stan.Subscription, error) {
	sc, err := stan.Connect(cfg.ClusterId, cfg.ClientId+"_sub")
	if err != nil {
		return nil, nil, err
	}

	sub, err := sc.Subscribe(cfg.Subject, cfg.subHandler)
	if err != nil {
		return nil, nil, err
	}

	if err := cfg.restoreCache(); err != nil {
		return nil, nil, err
	}
	return sc, sub, nil
}

// Handles incoming messages from the subject
func (c *StanConn) subHandler(m *stan.Msg) {
	var order model.Order
	var mut sync.Mutex

	mut.Lock()
	defer mut.Unlock()

	if err := json.Unmarshal(m.Data, &order); err != nil {
		logrus.Errorf("Could not read incoming message (%s)\n", err.Error())
		return
	}
	Msgs = append(Msgs, order)

	tx, err := c.DB.Begin()
	if err != nil {
		logrus.Errorf("Could not begin transaction (%s)\n", err.Error())
		return
	}

	query := fmt.Sprintf(`
	INSERT INTO %s
	(name, phone, zip, city, address, region, email)
	VALUES
	($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`, deliveriesTable)

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
		logrus.Errorf("Could not insert into deliveries table (%s)\n", err.Error())
		tx.Rollback()
		return
	}

	for _, item := range order.Items {
		query = fmt.Sprintf(`
		INSERT INTO %s
		(track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, chrt_id)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`, itemsTable)

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
			logrus.Errorf("Could not insert into items table (%s)\n", err.Error())
			tx.Rollback()
			return
		}

		order.ItemsIds = append(order.ItemsIds, int64(item.Id))
	}

	query = fmt.Sprintf(`
	INSERT INTO %s
	(transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`, paymentsTable)

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
		logrus.Errorf("Could not insert into payments table (%s)\n", err.Error())
		tx.Rollback()
		return
	}

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
		logrus.Errorf("Could not insert into orders table (%s)\n", err.Error())
		tx.Rollback()
		return
	}
	logrus.Print("Inserted order into the table")

	if err := tx.Commit(); err != nil {
		logrus.Errorf("Could not finish the transaction (%s)\n", err.Error())
	}
}

func (c *StanConn) restoreCache() error {
	var orders []model.Order
	logrus.Print("Test printing")

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

	if err := c.DB.Select(&orders, query); err != nil {
		return err
	}

	for _, order := range orders {
		for _, id := range order.ItemsIds {
			var item model.Item
			query = fmt.Sprintf(`
			SELECT * FROM %s
			WHERE id=$1
			`, itemsTable)

			if err := c.DB.Get(&item, query, id); err != nil {
				return err
			}

			order.Items = append(order.Items, item)
		}
	}

	Msgs = append(Msgs, orders...)

	return nil
}
