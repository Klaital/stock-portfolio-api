package nasdaq

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/klaital/stock-portfolio-api/stockfetcher"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Fetcher struct {
	URL     string
	Dataset string
	Key     string

	cacheLock         sync.RWMutex
	currentPriceCache map[string]decimal.Decimal
}

func New(apiKey string) *Fetcher {
	f := Fetcher{
		URL:               "https://data.nasdaq.com/api/v3/datasets/%s/%s.json?api_key=%s",
		Dataset:           "WIKI",
		Key:               apiKey,
		currentPriceCache: nil,
	}
	f.currentPriceCache = make(map[string]decimal.Decimal, 0)

	return &f
}

func (fetcher *Fetcher) computeURL(sym string) string {
	return fmt.Sprintf(fetcher.URL, fetcher.Dataset, sym, fetcher.Key)
}

func (fetcher *Fetcher) GetStockPrice(symbol string) (*stockfetcher.StockPrice, *time.Time, error) {
	// TODO: add caching

	resp, err := http.Get(fetcher.computeURL(symbol))
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"Status": resp.Status,
			"sym":    symbol,
			"url":    fetcher.computeURL(symbol),
		}).Fatal("Failed to fetch fresh data from price source")
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	nasdaqData := timeSeriesAndMetadata{}
	err = json.Unmarshal(respBody, &nasdaqData)
	if err != nil {
		log.WithError(err).Error("Failed to unmarshal NASDAQ response")
		return nil, nil, err
	}

	endDate, err := time.Parse("2006-01-02", nasdaqData.Dataset.EndDate)
	if err != nil {
		log.WithError(err).Fatal("Failed to parse response End Date")
	}
	yesterdayDate := "2000-01-01"
	stockPrice := stockfetcher.StockPrice{}
	foundToday := false

	for _, d := range nasdaqData.Dataset.Data {
		dayData := parseStockDayData(d)
		if dayData.Date == nasdaqData.Dataset.EndDate {
			stockPrice.Today = dayData.Close
			foundToday = true
		} else if strings.Compare(yesterdayDate, dayData.Date) < 0 {
			yesterdayDate = dayData.Date
			stockPrice.Yesterday = dayData.Close
		}
	}

	if !foundToday {
		return nil, nil, errors.New("data mismatch: latest data point not found in body")
	}
	return &stockPrice, &endDate, nil
}
