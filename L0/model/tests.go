package model

import "github.com/lib/pq"

type Test struct {
	Id    int           `json:"-" db:"id"`
	Value int           `json:"value" db:"value"`
	Text  string        `json:"text" db:"text"`
	Arr   pq.Int64Array `json:"arr" db:"arr"`
}
