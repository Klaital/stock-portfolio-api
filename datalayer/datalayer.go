package datalayer

import "time"

type StockStore interface {
	// Users

	AddUser(email, password string) error
	GetUserByEmail(email string) (*User, error)

	// Stock Positions
	GetPositionsByUser(userId uint64) ([]Position, error)
	AddPosition(userId uint64, symbol string, qty float64, basis uint64, boughtAt *time.Time) error
}
