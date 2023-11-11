package repository

import (
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

// Stan Connection configuration struct
// containing client_id`, `cluster_id`
// and `subject` to subscribe to
type StanConn struct {
	ClientId  string
	ClusterId string
	Subject   string
}

// Creates a new connection with Stan using info provided in the StanConn struct
// , returning connection interface, subscription interface and error
func NewStanConn(cfg StanConn) (stan.Conn, stan.Subscription, error) {
	sc, err := stan.Connect(cfg.ClusterId, cfg.ClientId+"_sub")
	if err != nil {
		return nil, nil, err
	}

	sub, err := sc.Subscribe(cfg.Subject, subHandler)
	if err != nil {
		return nil, nil, err
	}

	logrus.Debug("Subscribed to channel main")
	return sc, sub, nil
}

// Handles incoming messages from the subject
func subHandler(m *stan.Msg) {
	msg := string(m.Data)
	_ = msg
	// TODO: Handle msg
}
