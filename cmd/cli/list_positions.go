package main

import (
	"fmt"
	"github.com/klaital/stock-portfolio-api/datalayer"
	"github.com/klaital/stock-portfolio-api/stockfetcher"
	"github.com/klaital/stock-portfolio-api/stockfetcher/nasdaq"
	log "github.com/sirupsen/logrus"
)

func listPositions(datastore datalayer.StockStore, userId uint64) {
	positions, err := datastore.GetPositionsByUser(userId)
	if err != nil {
		log.WithField("userID", userId).WithError(err).Fatal("Failed to get positions for user")
	}

	var fetcher stockfetcher.StockFetcher
	fetcher = nasdaq.New()

	fmt.Print("Sym\tQty\tBasis\tCurrent\n")
	for _, p := range positions {
		closePrice, _, err := fetcher.GetStockPrice(p.Symbol)
		if err != nil {
			log.WithField("sym", p.Symbol).WithError(err).Fatal("Failed to fetch stock price")
		}
		fmt.Printf("%s\t%s\t%s\t%s\n", p.Symbol, p.Qty.String(), p.Basis.String(), closePrice.String())
	}
}
