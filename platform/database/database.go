package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewDB creates a new GORM database connection
func NewDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
