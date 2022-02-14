package nasdaq

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sync"
	"time"
)

type Fetcher struct {
	URL     string
	Dataset string

	cacheLock         sync.RWMutex
	currentPriceCache map[string]decimal.Decimal
}

func New() *Fetcher {
	f := Fetcher{
		URL:               "https://data.nasdaq.com/api/v3/datasets/%s/%s.json",
		Dataset:           "WIKI",
		currentPriceCache: nil,
	}
	f.currentPriceCache = make(map[string]decimal.Decimal, 0)

	return &f
}

func (fetcher *Fetcher) computeURL(sym string) string {
	return fmt.Sprintf(fetcher.URL, fetcher.Dataset, sym)
}

func (fetcher *Fetcher) GetStockPrice(symbol string) (*decimal.Decimal, *time.Time, error) {
	// TODO: add caching
	
	resp, err := http.Get(fetcher.computeURL(symbol))
	if err != nil {
		return nil, nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	nasdaqData := timeSeriesAndMetadata{}
	err = json.Unmarshal(respBody, &nasdaqData)
	if err != nil {
		log.WithError(err).Error("Failed to unmarshal NASDAQ response")
		return nil, nil, err
	}

	endTime, err := time.Parse("2006-01-02", nasdaqData.Dataset.EndDate)
	if err != nil {
		log.WithField("EndDate", nasdaqData.Dataset.EndDate).WithError(err).Error("Failed to parse end date")
		return nil, nil, err
	}

	for _, d := range nasdaqData.Dataset.Data {
		dayData := parseStockDayData(d)
		if dayData.Date == nasdaqData.Dataset.EndDate {
			log.Debug("Found latest data point")
			return &dayData.Close, &endTime, nil
		}
	}

	return nil, nil, errors.New("data mismatch: latest data point not found in body")
}
