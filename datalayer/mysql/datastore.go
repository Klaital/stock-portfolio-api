package mysql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type DataStore struct {
	db       *sql.DB
	ctx      context.Context
	HashCost int
}

func New(ctx context.Context, host, username, password, dbName string, port, hashCost int) (*DataStore, error) {
	store := &DataStore{
		db:       nil,
		ctx:      ctx,
		HashCost: hashCost,
	}
	var err error
	store.db, err = sql.Open("mysql", dsn(host, username, password, dbName, port))
	if err != nil {
		log.WithError(err).Error("Failed to connect to DB")
		return nil, err
	}
	return store, nil
}

func dsn(host, username, password, dbName string, port int) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", username, password, host, port, dbName)
}
