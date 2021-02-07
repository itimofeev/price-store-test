package model

import "time"

type ParsedProduct struct {
	Name  string
	Price int64
}

type Product struct {
	ID          int64
	Name        string
	Price       int64
	LastUpdate  time.Time
	UpdateCount int64
}
