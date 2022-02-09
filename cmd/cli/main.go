package main

import (
	"context"
	"flag"
	"github.com/klaital/stock-portfolio-api/datalayer/mysql"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func main() {
	addNewPositionCmd := flag.NewFlagSet("add-stock", flag.ExitOnError)
	addParams := struct {
		userId   uint64
		symbol   string
		qty      float64
		boughtAt string
		basis    uint64
	}{}
	addNewPositionCmd.StringVar(&addParams.symbol, "sym", "", "Stock Symbol")
	addNewPositionCmd.Float64Var(&addParams.qty, "qty", 0.0, "Quantity bought")
	addNewPositionCmd.StringVar(&addParams.symbol, "at", "", "When was the stock bought. Use ISO8601 - YYYY-MM-DD")
	addNewPositionCmd.Uint64Var(&addParams.userId, "user", 0, "User ID")

	datastore, err := mysql.New(context.Background(), "mysql.abandonedfactory.net", "klaital", "h0Shinokoe", "stocks", 3306, bcrypt.MaxCost)
	if err != nil {
		log.WithError(err).Fatal("Failed to init datastore")
	}

	switch os.Args[1] {
	case "add-stock":
		addNewPositionCmd.Parse(os.Args[2:])

		addPosition(datastore, addParams.userId, addParams.symbol, addParams.qty, addParams.basis, addParams.boughtAt)
	}
}
