package datalayer

import (
	"github.com/shopspring/decimal"
	"time"
)

type StockStore interface {
	// Users

	AddUser(email, password string) error
	GetUserByEmail(email string) (*User, error)

	// Stock Positions
	GetPositionsByUser(userId uint64) ([]Position, error)
	AddPosition(userId uint64, symbol string, qty decimal.Decimal, basis decimal.Decimal, boughtAt time.Time) error
}
