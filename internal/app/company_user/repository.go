package company_user

import (
	"crypto/sha256"
	"fmt"
	"time"

	"gobizmanager/internal/app/pkg/model"
	"gobizmanager/internal/app/rbac"
	"gobizmanager/internal/app/user"
	"gobizmanager/pkg/encryption"

	"gobizmanager/platform/config"

	"gorm.io/gorm"
)

type Repository struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewRepository(db *gorm.DB, cfg *config.Config) *Repository {
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

	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
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
	hashedPassword, err := encryption.HashPassword(req.Password)
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
	var userRole rbac.Role
	if err := tx.Where("name = ? AND company_id IS NULL", "USER").First(&userRole).Error; err != nil {
		return nil, fmt.Errorf("failed to get USER role ID: %w", err)
	}

	// Assign USER role to the new user
	userRoleAssignment := &rbac.UserRole{
		UserID:    userID,
		RoleID:    userRole.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := tx.Create(userRoleAssignment).Error; err != nil {
		return nil, fmt.Errorf("failed to assign USER role: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return companyUser, nil
}

func (r *Repository) createUser(tx *gorm.DB, username, emailHash, password, phone string) (int64, error) {
	now := time.Now()
	user := &model.User{
		Email:     username,
		EmailHash: emailHash,
		Password:  password,
		Phone:     phone,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := tx.Create(user).Error; err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return user.ID, nil
}

func (r *Repository) createCompanyUser(tx *gorm.DB, companyID, userID int64) (*CompanyUser, error) {
	now := time.Now()
	companyUser := &CompanyUser{
		CompanyID: companyID,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := tx.Create(companyUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create company-user relationship: %w", err)
	}

	return companyUser, nil
}

func (r *Repository) ListCompanyUsers(companyID int64) ([]*CompanyUser, error) {
	var users []*CompanyUser
	if err := r.db.Where("company_id = ?", companyID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) RemoveCompanyUser(companyID, userID int64) error {
	return r.db.Where("company_id = ? AND user_id = ?", companyID, userID).Delete(&CompanyUser{}).Error
}
