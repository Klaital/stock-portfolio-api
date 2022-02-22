package main

import (
	"flag"
	"fmt"
	"github.com/klaital/stock-portfolio-api/datalayer"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func addPosition(datastore datalayer.StockStore) {
	var err error

	cmd := flag.NewFlagSet("add-stock", flag.ExitOnError)
	params := struct {
		userId   uint64
		symbol   string
		qty      string
		boughtAt string
		basis    string
	}{}
	cmd.StringVar(&params.symbol, "sym", "", "Stock Symbol")
	cmd.StringVar(&params.qty, "qty", "", "Quantity bought")
	cmd.StringVar(&params.boughtAt, "at", "", "When was the stock bought. Use ISO8601 - YYYY-MM-DD")
	cmd.Uint64Var(&params.userId, "user", 0, "User ID")
	cmd.StringVar(&params.basis, "basis", "", "User ID")
	if err = cmd.Parse(os.Args[2:]); err != nil {
		log.WithError(err).Fatal("Failed to parse commandline params")
	}
	log.WithFields(log.Fields{
		"symbol":   params.symbol,
		"userID":   params.userId,
		"qty":      params.qty,
		"basis":    params.basis,
		"boughtAt": params.boughtAt,
	}).Debug("Got CLI params")

	buyTimestamp := time.Now()
	if params.boughtAt != "" {
		buyTimestamp, err = time.Parse("2006-01-02", params.boughtAt)
		if err != nil {
			log.WithField("raw", params.boughtAt).WithError(err).Error("Failed to parse timestamp. Using NOW instead.")
			buyTimestamp = time.Now()
		}
	}

	// Convert the floating-point strings into fixed-point decimals
	var qty, basis decimal.Decimal
	if qty, err = decimal.NewFromString(params.qty); err != nil {
		log.WithField("raw", params.qty).Fatal("Failed to parse qty")
	}
	if basis, err = decimal.NewFromString(params.basis); err != nil {
		log.WithField("raw", params.basis).Fatal("Failed to parse basis")
	}
	if err = datastore.AddPosition(params.userId, params.symbol, qty, basis, buyTimestamp); err != nil {
		log.WithError(err).WithField("params", params).Fatal("Error from datalayer when adding a new position")
	}
	fmt.Printf("Added %s shares of %s\n", qty.String(), params.symbol)
}
