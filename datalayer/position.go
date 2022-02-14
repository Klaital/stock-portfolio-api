package datalayer

import (
	"github.com/shopspring/decimal"
	"time"
)

type Position struct {
	ID        uint64
	UserID    uint64
	Symbol    string
	Qty       decimal.Decimal
	Basis     decimal.Decimal
	BoughtAt  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p Position) Value(priceEach decimal.Decimal) decimal.Decimal {
	return p.Qty.Mul(priceEach)
}
func (p Position) BasisTotal() decimal.Decimal {
	return p.Qty.Mul(p.Basis)
}

func (p Position) GainAmount(priceEach decimal.Decimal) decimal.Decimal {
	return p.Value(priceEach).Sub(p.BasisTotal())
}
