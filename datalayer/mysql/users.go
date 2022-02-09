package mysql

import (
	"database/sql"
	"github.com/klaital/stock-portfolio-api/datalayer"
	log "github.com/sirupsen/logrus"
)

// AddUser inserts a new user record by hashing the given password.
func (store *DataStore) AddUser(email, password string) error {
	passwordDigest, err := datalayer.HashAndSalt(password, store.HashCost)
	if err != nil {
		log.WithError(err).Error("Failed to hash password")
		return err
	}
	_, err = store.db.ExecContext(store.ctx, `INSERT INTO users (email, password_digest) VALUES (?, ?)`, email, passwordDigest)
	if err != nil {
		log.WithError(err).Error("Failed to insert new user record")
		return err
	}

	// Success!
	return nil
}

// GetUserByEmail fetches a user record from the DB
func (store *DataStore) GetUserByEmail(email string) (*datalayer.User, error) {
	rows, err := store.db.Query(`SELECT user_id, created_at, updated_at, password_digest FROM users WHERE email = ?`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()
	u := datalayer.User{
		Email: email,
	}
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.PasswordDigest)
		if err != nil {
			log.WithError(err).Error("Error scanning user record")
			return nil, err
		}
	}
	return &u, nil
}
