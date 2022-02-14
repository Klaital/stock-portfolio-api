package main

import (
	"fmt"
	"github.com/klaital/stock-portfolio-api/datalayer"
	log "github.com/sirupsen/logrus"
)

func addUser(datastore datalayer.StockStore, email, password string) {
	err := datastore.AddUser(email, password)
	if err != nil {
		log.WithError(err).Fatal("Failed to add new user")
	}

	u, err := datastore.GetUserByEmail(email)
	if err != nil {
		log.WithError(err).Fatal("Failed to fetch back new user")
	}
	fmt.Printf("%d\t%s\n", u.ID, u.Email)
}
