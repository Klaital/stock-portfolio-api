package main

import (
	"github.com/klaital/stock-portfolio-api/datalayer"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func addPosition(datastore datalayer.StockStore, userId uint64, symbol string, qty float64, basis uint64, boughtAt string) {
	buyTimestamp := time.Now()
	var err error
	if boughtAt != "" {
		buyTimestamp, err = time.Parse(time.RFC3339, boughtAt)
		if err != nil {
			log.WithField("raw", boughtAt).WithError(err).Error("Failed to parse timestamp. Using NOW instead.")
			buyTimestamp = time.Now()
		}
	}
	err = datastore.AddPosition(userId, symbol, qty, basis, &buyTimestamp)
	if err != nil {
		log.WithError(err).Error("Error from datalayer when adding a new position")
		os.Exit(1)
	}
}
