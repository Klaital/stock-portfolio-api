package mysql

import (
	"github.com/klaital/stock-portfolio-api/datalayer"
	log "github.com/sirupsen/logrus"
	"time"
)

func (store *DataStore) AddPosition(userId uint64, symbol string, qty float64, basis uint64, boughtAt *time.Time) error {
	_, err := store.db.ExecContext(store.ctx, `INSERT INTO positions (user_id, sytmbol, qty, basis, bought_at) VALUES (?, ?, ?, ?, ?)`, userId, symbol, qty, basis, boughtAt)
	return err
}

func (store *DataStore) GetPositionsBySymbol(userId uint64, symbol string) ([]datalayer.Position, error) {
	rows, err := store.db.QueryContext(store.ctx, `SELECT position_id, qty, basis, bought_at, created_at, updated_at FROM positions WHERE user_id = ? AND symbol = ?`, userId, symbol)
	if err != nil {
		return nil, err
	}
	positions := make([]datalayer.Position, 0)
	var tmp *datalayer.Position
	defer rows.Close()
	for rows.Next() {
		tmp = new(datalayer.Position)
		err = rows.Scan(&tmp.ID, &tmp.Qty, &tmp.Basis, &tmp.BoughtAt, &tmp.CreatedAt, &tmp.UpdatedAt)
		if err != nil {
			log.WithError(err).Error("Error scanning position")
			return nil, err
		}
		tmp.UserID = userId
		tmp.Symbol = symbol
		positions = append(positions, *tmp)
	}

	return positions, nil
}