package company_user

import (
	"database/sql"
	"fmt"
	"time"

	"gobizmanager/internal/user"
	"gobizmanager/pkg/encryption"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// RegisterCompanyUser registers a new user for a company
func (r *Repository) RegisterCompanyUser(req *RegisterCompanyUserRequest) (*CompanyUser, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Encrypt sensitive data
	encryptedUsername, err := encryption.Encrypt(req.Username, "username")
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt username: %w", err)
	}

	encryptedPhone, err := encryption.Encrypt(req.Phone, "phone")
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt phone: %w", err)
	}

	// Hash password
	hashedPassword, err := user.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	userID, err := r.createUser(tx, encryptedUsername, hashedPassword, encryptedPhone)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create company-user relationship
	companyUser, err := r.createCompanyUser(tx, req.CompanyID, userID, req.IsMain)
	if err != nil {
		return nil, fmt.Errorf("failed to create company-user relationship: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return companyUser, nil
}

func (r *Repository) createUser(tx *sql.Tx, username, password, phone string) (int64, error) {
	result, err := tx.Exec(`
		INSERT INTO users (username, password, phone, created_at, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, username, password, phone)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return result.LastInsertId()
}

func (r *Repository) createCompanyUser(tx *sql.Tx, companyID, userID int64, isMain bool) (*CompanyUser, error) {
	now := time.Now()
	result, err := tx.Exec(`
		INSERT INTO company_users (company_id, user_id, is_main, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, companyID, userID, isMain, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create company-user relationship: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get company-user ID: %w", err)
	}

	return &CompanyUser{
		ID:        id,
		CompanyID: companyID,
		UserID:    userID,
		IsMain:    isMain,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
