package model

// Secondary item object
type Item struct {
	Id int `json:"id"`
	TrackNumber string `json:"track_number"`
	Price int `json:"price"`
	RId string `json:"rid"`
	Name string `json:"name"`
	Sale int `json:"sale"`
	Size string `json:"size"`
	TotalPrice int `json:"total_price"`
	NmId int `json:"nm_id"`
	Brand string `json:"brand"`
	Status int `json:"status"`
	ChartId int `json:"chrt_id"`
}