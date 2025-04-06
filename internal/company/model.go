package company

import (
	"time"

	"gobizmanager/pkg/encryption"
)

type Company struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone" encrypted:"true"`
	Email       string    `json:"email" encrypted:"true"`
	Identifier  string    `json:"identifier"`
	Logo        string    `json:"logo"`
	RootGroupID int64     `json:"root_group_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EncryptSensitiveFields encrypts sensitive fields using the provided key
func (c *Company) EncryptSensitiveFields(key string) error {
	var err error
	if c.Phone != "" {
		c.Phone, err = encryption.Encrypt(c.Phone, key)
		if err != nil {
			return err
		}
	}
	if c.Email != "" {
		c.Email, err = encryption.Encrypt(c.Email, key)
		if err != nil {
			return err
		}
	}
	return nil
}

// DecryptSensitiveFields decrypts sensitive fields using the provided key
func (c *Company) DecryptSensitiveFields(key string) error {
	var err error
	if c.Phone != "" {
		c.Phone, err = encryption.Decrypt(c.Phone, key)
		if err != nil {
			return err
		}
	}
	if c.Email != "" {
		c.Email, err = encryption.Decrypt(c.Email, key)
		if err != nil {
			return err
		}
	}
	return nil
}
