package model

// Secondary payment object
type Payment struct {
	Id int `json:"id"`
	Transaction string `json:"transaction"`
	RequestID string `json:"request_id"`
	Currency string `json:"currency"`
	Provider string `json:"provider"`
	Amount int `json:"amount"`
	PaymentDT int `json:"payment_dt"`
	Bank string `json:"bank"`
	DeliveryCost int `json:"delivery_cost"`
	GoodsTotal int `json:"goods_total"`
	CustomFee int `json:"custom_fee"`
}