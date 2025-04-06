package company

import (
	"database/sql"
	"fmt"
	"time"

	"gobizmanager/pkg/encryption"
)

type Company struct {
	ID         int64          `json:"id"`
	Name       string         `json:"name"`
	Phone      string         `json:"phone" encrypted:"true"`
	Email      string         `json:"email" encrypted:"true"`
	Address    string         `json:"address" encrypted:"true"`
	Identifier string         `json:"identifier"`
	Logo       sql.NullString `json:"logo"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// EncryptSensitiveFields encrypts sensitive fields using the provided key
func (c *Company) EncryptSensitiveFields(key string) error {
	if c.Phone != "" {
		encrypted, err := encryption.Encrypt(c.Phone, key)
		if err != nil {
			return fmt.Errorf("failed to encrypt phone: %w", err)
		}
		c.Phone = encrypted
	}

	if c.Email != "" {
		encrypted, err := encryption.Encrypt(c.Email, key)
		if err != nil {
			return fmt.Errorf("failed to encrypt email: %w", err)
		}
		c.Email = encrypted
	}

	if c.Address != "" {
		encrypted, err := encryption.Encrypt(c.Address, key)
		if err != nil {
			return fmt.Errorf("failed to encrypt address: %w", err)
		}
		c.Address = encrypted
	}

	return nil
}

// DecryptSensitiveFields decrypts sensitive fields using the provided key
func (c *Company) DecryptSensitiveFields(key string) error {
	if c.Phone != "" {
		decrypted, err := encryption.Decrypt(c.Phone, key)
		if err != nil {
			return fmt.Errorf("failed to decrypt phone: %w", err)
		}
		c.Phone = decrypted
	}

	if c.Email != "" {
		decrypted, err := encryption.Decrypt(c.Email, key)
		if err != nil {
			return fmt.Errorf("failed to decrypt email: %w", err)
		}
		c.Email = decrypted
	}

	if c.Address != "" {
		decrypted, err := encryption.Decrypt(c.Address, key)
		if err != nil {
			return fmt.Errorf("failed to decrypt address: %w", err)
		}
		c.Address = decrypted
	}

	return nil
}
