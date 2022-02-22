package stockfetcher

import (
	"time"
)

type StockFetcher interface {
	GetStockPrice(symbol string) (*StockPrice, *time.Time, error)
}
