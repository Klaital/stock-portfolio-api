package main

import (
	"fmt"
	"github.com/klaital/stock-portfolio-api/datalayer"
	"github.com/klaital/stock-portfolio-api/stockfetcher"
	fhub "github.com/klaital/stock-portfolio-api/stockfetcher/finnhub"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func listPositions(datastore datalayer.StockStore, userId uint64) {
	positions, err := datastore.GetPositionsByUser(userId)
	if err != nil {
		log.WithField("userID", userId).WithError(err).Fatal("Failed to get positions for user")
	}

	var fetcher stockfetcher.StockFetcher
	fetcher = fhub.New(viper.GetString("FINNHUB_API_KEY"))

	fmt.Print("Sym\tQty\tBasis\tCurrent\tOverall\n")
	total := stockfetcher.StockPrice{}
	for _, p := range positions {
		closePrice, _, err := fetcher.GetStockPrice(p.Symbol)
		if err != nil {
			log.WithField("sym", p.Symbol).WithError(err).Fatal("Failed to fetch stock price")
		}
		total.Today = total.Today.Add(closePrice.Today)
		total.Yesterday = total.Yesterday.Add(closePrice.Yesterday)

		fmt.Printf("%s\t%s\t%s\t%s\t%s\n", p.Symbol, p.Qty.String(), p.Basis.String(), closePrice.TodayWithChange(), closePrice.TodayWithChangeByQty(p.Qty))
	}

	fmt.Printf("\n\nOverall Current Value: %s\n", total.TodayWithChange())
}
