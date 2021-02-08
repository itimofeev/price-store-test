package model

import "time"

type ParsedProduct struct {
	Name  string
	Price int64
}

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Price       int64     `json:"price"`
	LastUpdate  time.Time `json:"lastUpdate"`
	UpdateCount int64     `json:"updateCount"`
}
