package main

import (
	"fmt"
	"github.com/klaital/stock-portfolio-api/datalayer"
	log "github.com/sirupsen/logrus"
)

func listPositions(datastore datalayer.StockStore, userId uint64) {
	positions, err := datastore.GetPositionsByUser(userId)
	if err != nil {
		log.WithField("userID", userId).WithError(err).Fatal("Failed to get positions for user")
	}

	for _, p := range positions {
		fmt.Printf("%s\t%f\n%d\n", p.Symbol, p.Qty, p.Basis)
	}
}
