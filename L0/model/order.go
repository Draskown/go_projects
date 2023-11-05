package model

import "time"

// Nested order structure.
// Contains info about another
// JSON objects
type Order struct {
	OrderId string `json:"order_uid"`
	TrackNumber string `json:"track_number"`
	Entry string `json:"entry"`
	Delivery Delivery `json:"delivery"`
	Payment Payment `json:"payment"`
	Items []Item `json:"items"`
	Locale string `json:"locale"`
	InternalSignature string `json:"internal_signature"`
	CustomerId string `json:"customer_id"`
	DeliveryService string `json:"delivery_service"`
	ShardKey string `json:"shardkey"`
	SmId int `json:"sm_id"`
	DateCreated string `json:"date_created"`
	OofShard string `json:"oof_shard"`
}

// Converts string date from the struct to time format
func (o *Order) GetDate () (time.Time, error) {
	layout := "2006-01-02T15:04:05Z"

	return time.Parse(layout, o.DateCreated)
}

// Sets the date from time format as string into the structure
func (o *Order) SetDate (in time.Time) error {
	temp, err := in.UTC().MarshalText()
	if err != nil{
		return err
	}

	o.DateCreated = string(temp)
	return nil
}
