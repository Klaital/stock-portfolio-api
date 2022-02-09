package datalayer

import "time"

type Position struct {
	ID        uint64
	UserID    uint64
	Symbol    string
	Qty       float64
	Basis     uint64
	BoughtAt  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
