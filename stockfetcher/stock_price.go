package stockfetcher

import (
	"fmt"
	"github.com/shopspring/decimal"
)

type StockPrice struct {
	Today     decimal.Decimal
	Yesterday decimal.Decimal
}

func (s StockPrice) Delta() decimal.Decimal {
	return s.Today.Sub(s.Yesterday)
}

func (s StockPrice) TodayWithChange() string {
	prefix := "+"
	if s.Yesterday.GreaterThan(s.Today) {
		prefix = ""
	}

	return fmt.Sprintf("%s (%s%s)", s.Today.Round(2).String(), prefix, s.Today.Sub(s.Yesterday).Round(2))
}

func (s StockPrice) TodayWithChangeByQty(qty decimal.Decimal) string {
	prefix := "+"
	if s.Yesterday.GreaterThan(s.Today) {
		prefix = ""
	}

	return fmt.Sprintf("%s (%s%s)", s.Today.Mul(qty).Round(2), prefix, s.Delta().Mul(qty).Round(2))
}
