package repository

import (
	"encoding/json"
	"sync"

	"github.com/Draskown/WBL0/model"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

// Cache slice
var Msgs []model.Order

// Stan Connection configuration struct
// containing client_id`, `cluster_id`
// and `subject` to subscribe to
type StanCfg struct {
	ClientId  string
	ClusterId string
	Subject   string
}

// Creates a new connection with Stan using info provided in the StanConn struct
// , returning connection interface, subscription interface and error
func (r *DBConvPostgres) ConnectStan(cfg StanCfg) (stan.Conn, stan.Subscription, error) {
	// Restore cache from DB
	if err := r.restoreCache(); err != nil {
		return nil, nil, err
	}

	// Connect to STAN cluster
	sc, err := stan.Connect(cfg.ClusterId, cfg.ClientId+"_sub")
	if err != nil {
		return nil, nil, err
	}

	// Subscribe to the subject
	sub, err := sc.Subscribe(cfg.Subject, r.subHandler)
	if err != nil {
		return nil, nil, err
	}

	return sc, sub, nil
}

// Handles incoming messages from the subject
//
// Adds the message to DB and Msgs cache
func (r *DBConvPostgres) subHandler(m *stan.Msg) {
	var order model.Order
	var mut sync.Mutex

	// Lock operated structs and cache array
	mut.Lock()
	defer mut.Unlock()

	// Get order from JSON
	if err := json.Unmarshal(m.Data, &order); err != nil {
		logrus.Errorf("Could not read incoming message (%s)\n", err.Error())
		return
	}
	// Append order to cache
	Msgs = append(Msgs, order)

	if err := r.insertOrder(order); err != nil {
		logrus.Error(err.Error())
	}
}
