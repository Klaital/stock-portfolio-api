package stockfetcher

import (
	"github.com/shopspring/decimal"
	"time"
)

type StockFetcher interface {
	GetStockPrice(symbol string) (*decimal.Decimal, *time.Time, error)
}
