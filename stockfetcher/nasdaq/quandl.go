package nasdaq

import (
	"github.com/shopspring/decimal"
	"time"
)

type timeSeriesAndMetadata struct {
	Dataset struct {
		ID                  int             `json:"id"`
		DatasetCode         string          `json:"dataset_code"`
		DatabaseCode        string          `json:"database_code"`
		Name                string          `json:"name"`
		Description         string          `json:"description"`
		RefreshedAt         time.Time       `json:"refreshed_at"`
		NewestAvailableDate string          `json:"newest_available_date"`
		OldestAvailableDate string          `json:"oldest_available_date"`
		ColumnNames         []string        `json:"column_names"`
		Frequency           string          `json:"frequency"`
		Type                string          `json:"type"`
		Premium             bool            `json:"premium"`
		StartDate           string          `json:"start_date"`
		EndDate             string          `json:"end_date"`
		Data                [][]interface{} `json:"data"`
		Collapse            interface{}     `json:"collapse"`
		Order               string          `json:"order"`
		DatabaseID          int             `json:"database_id"`
	} `json:"dataset"`
}

type StockDay struct {
	Date   string
	Open   decimal.Decimal
	High   decimal.Decimal
	Low    decimal.Decimal
	Close  decimal.Decimal
	Volume decimal.Decimal
}

func parseStockDayData(data []interface{}) StockDay {
	s := StockDay{}
	s.Date = data[0].(string)
	s.Open = decimal.NewFromFloat(data[1].(float64))
	s.High = decimal.NewFromFloat(data[2].(float64))
	s.Low = decimal.NewFromFloat(data[3].(float64))
	s.Close = decimal.NewFromFloat(data[4].(float64))
	s.Volume = decimal.NewFromFloat(data[5].(float64))
	return s
}
