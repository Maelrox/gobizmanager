package company_user

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"

	"gobizmanager/internal/user"
	"gobizmanager/pkg/encryption"
	"gobizmanager/platform/config"
)

type Repository struct {
	db  *sql.DB
	cfg *config.Config
}

func NewRepository(db *sql.DB, cfg *config.Config) *Repository {
	return &Repository{db: db, cfg: cfg}
}

// RegisterCompanyUser registers a new user for a company
func (r *Repository) RegisterCompanyUser(req *RegisterCompanyUserRequest) (*CompanyUser, error) {
	// Check if user already exists
	userRepo := user.NewRepository(r.db, r.cfg)
	_, err := userRepo.GetUserByEmail(req.Username)
	if err == nil {
		return nil, fmt.Errorf("username already exists")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Encrypt sensitive data
	encryptedUsername, err := encryption.Encrypt(req.Username, r.cfg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt username: %w", err)
	}

	encryptedPhone, err := encryption.Encrypt(req.Phone, r.cfg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt phone: %w", err)
	}

	// Hash password
	hashedPassword, err := user.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create email hash for searching
	emailHash := fmt.Sprintf("%x", sha256.Sum256([]byte(req.Username)))

	// Create user
	userID, err := r.createUser(tx, encryptedUsername, emailHash, hashedPassword, encryptedPhone)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create company-user relationship
	companyUser, err := r.createCompanyUser(tx, req.CompanyID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create company-user relationship: %w", err)
	}

	// Get USER role ID
	var userRoleID int64
	err = tx.QueryRow("SELECT id FROM roles WHERE name = 'USER' AND company_id IS NULL").Scan(&userRoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get USER role ID: %w", err)
	}

	// Assign USER role to the new user
	_, err = tx.Exec(`
		INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, userID, userRoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to assign USER role: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return companyUser, nil
}

func (r *Repository) createUser(tx *sql.Tx, username, emailHash, password, phone string) (int64, error) {
	result, err := tx.Exec(`
		INSERT INTO users (email, email_hash, password, phone, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, username, emailHash, password, phone)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return result.LastInsertId()
}

func (r *Repository) createCompanyUser(tx *sql.Tx, companyID, userID int64) (*CompanyUser, error) {
	now := time.Now()
	result, err := tx.Exec(`
		INSERT INTO company_users (company_id, user_id, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, companyID, userID, now, now)
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
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
