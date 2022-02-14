package datalayer

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID             uint64
	Email          string
	PasswordDigest string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func HashAndSalt(pwd string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password))
	return err == nil
}
