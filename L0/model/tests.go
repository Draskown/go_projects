package model

import "github.com/lib/pq"

type Test struct {
	Value int `json:"value"`
	Text string `json:"text"`
	Arr pq.Int64Array `json:"arr"`
	Arr_One pq.Int64Array `json:"arr_one"`
}