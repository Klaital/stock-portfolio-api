package finnhub

import (
	"context"
	"fmt"
	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/klaital/stock-portfolio-api/stockfetcher"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Fetcher struct {
	client *finnhub.DefaultApiService

	cacheLock         sync.RWMutex
	currentPriceCache map[string]decimal.Decimal
}

func New(apiKey string) *Fetcher {
	f := Fetcher{}
	f.currentPriceCache = make(map[string]decimal.Decimal, 0)
	cfg := finnhub.NewConfiguration()
	if len(apiKey) == 0 {
		log.Fatal("Finnhub key missing")
	}
	cfg.AddDefaultHeader("X-Finnhub-Token", apiKey)
	f.client = finnhub.NewAPIClient(cfg).DefaultApi

	return &f
}

func (f *Fetcher) GetStockPrice(symbol string) (*stockfetcher.StockPrice, *time.Time, error) {
	quote, _, err := f.client.Quote(context.Background()).Symbol(symbol).Execute()
	if err != nil {
		return nil, nil, err
	}
	today, ok := quote.GetCOk()
	if !ok {
		return nil, nil, fmt.Errorf("no current price in stock quote for %s", symbol)
	}
	yesterday, ok := quote.GetPcOk()
	if !ok {
		return nil, nil, fmt.Errorf("no previous close price in stock quote for %s", symbol)
	}
	t := time.Now()
	return &stockfetcher.StockPrice{
		Today:     decimal.NewFromFloat32(*today),
		Yesterday: decimal.NewFromFloat32(*yesterday),
	}, &t, nil
}
