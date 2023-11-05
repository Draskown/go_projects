package model

// Secondary delivery object
type Delivery struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	ZipCode string `json:"zip"`
	City string `json:"city"`
	Address string `json:"address"`
	Region string `json:"region"`
	Email string `json:"email"`
}