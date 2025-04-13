package model

import (
	"time"

	utils "gobizmanager/pkg/encryption"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email" encrypted:"true"`
	EmailHash string    `json:"-" gorm:"index"`
	Password  string    `json:"-"`
	Phone     string    `json:"phone" encrypted:"true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) EncryptSensitiveFields(key string) error {
	var err error
	if u.Email != "" {
		u.Email, err = utils.Encrypt(u.Email, key)
		if err != nil {
			return err
		}
	}
	if u.Phone != "" {
		u.Phone, err = utils.Encrypt(u.Phone, key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) DecryptSensitiveFields(key string) error {
	var err error
	if u.Email != "" {
		u.Email, err = utils.Decrypt(u.Email, key)
		if err != nil {
			return err
		}
	}
	if u.Phone != "" {
		u.Phone, err = utils.Decrypt(u.Phone, key)
		if err != nil {
			return err
		}
	}
	return nil
}
