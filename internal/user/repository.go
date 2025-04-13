// internal/user/repository.go
package user

import (
	"crypto/sha256"
	"fmt"
	"time"

	model "gobizmanager/internal/models"
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

// CreateUserWithTx creates a new user within a transaction
func (r *Repository) CreateUserWithTx(tx *gorm.DB, username, password, phone string) (int64, error) {
	hashedPassword, err := encryption.HashPassword(password)
	if err != nil {
		return 0, err
	}

	// Create a temporary user to encrypt fields
	user := &model.User{
		Email:    username,
		Password: hashedPassword,
		Phone:    phone,
	}

	// Encrypt sensitive fields
	if err := user.EncryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return 0, err
	}

	// Create email hash for searching
	emailHash := fmt.Sprintf("%x", sha256.Sum256([]byte(username)))

	now := time.Now()
	user.EmailHash = emailHash
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := tx.Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

// GetUserByID returns a user by ID
func (r *Repository) GetUserByID(id int64) (*model.User, error) {
	user := &model.User{}
	if err := r.db.First(user, id).Error; err != nil {
		return nil, err
	}

	// Decrypt sensitive fields
	if err := user.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetUserByEmail(email string) (*model.User, error) {
	emailHash := fmt.Sprintf("%x", sha256.Sum256([]byte(email)))

	user := &model.User{}
	if err := r.db.Where("email_hash = ?", emailHash).First(user).Error; err != nil {
		return nil, err
	}

	// Decrypt sensitive fields
	if err := user.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return nil, err
	}

	return user, nil
}

// GetDB returns the underlying database connection
func (r *Repository) GetDB() *gorm.DB {
	return r.db
}

// CreateUser creates a new user
func (r *Repository) CreateUser(email, password, phone string) (int64, error) {
	hashedPassword, err := encryption.HashPassword(password)
	if err != nil {
		return 0, err
	}

	// Create a temporary user to encrypt fields
	user := &model.User{
		Email:    email,
		Password: hashedPassword,
		Phone:    phone,
	}

	// Encrypt sensitive fields
	if err := user.EncryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return 0, err
	}

	// Create email hash for searching
	emailHash := fmt.Sprintf("%x", sha256.Sum256([]byte(email)))
	user.EmailHash = emailHash

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := r.db.Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

// UpdateUser updates a user
func (r *Repository) UpdateUser(id int64, email, password, phone string) error {
	user := &model.User{}
	if err := r.db.First(user, id).Error; err != nil {
		return err
	}

	// Update fields
	if email != "" {
		user.Email = email
		user.EmailHash = fmt.Sprintf("%x", sha256.Sum256([]byte(email)))
	}
	if password != "" {
		hashedPassword, err := encryption.HashPassword(password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}
	if phone != "" {
		user.Phone = phone
	}

	// Encrypt sensitive fields
	if err := user.EncryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return err
	}

	user.UpdatedAt = time.Now()
	return r.db.Save(user).Error
}

// GetRootRoleID returns the root role ID
func (r *Repository) GetRootRoleID(tx *gorm.DB) (int64, error) {
	var roleID int64
	if err := tx.Model(&model.Role{}).Where("name = ?", "ROOT").Select("id").First(&roleID).Error; err != nil {
		return 0, err
	}
	return roleID, nil
}

// AssignRootRole assigns the root role to a user
func (r *Repository) AssignRootRole(tx *gorm.DB, userID, roleID int64) error {
	userRole := &model.UserRole{
		UserID:    userID,
		RoleID:    roleID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return tx.Create(userRole).Error
}

// RegisterRootUser registers a root user
func (r *Repository) RegisterRootUser(username, password string) (int64, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	// Create user
	userID, err := r.CreateUserWithTx(tx, username, password, "")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get root role ID
	roleID, err := r.GetRootRoleID(tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Assign root role
	if err := r.AssignRootRole(tx, userID, roleID); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return userID, nil
}

// RegisterUser registers a new user
func (r *Repository) RegisterUser(username, password, phone string) (int64, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	// Create user
	userID, err := r.CreateUserWithTx(tx, username, password, phone)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *Repository) IsRoot(userID int64) (bool, error) {
	var count int64
	if err := r.db.Model(&model.UserRole{}).
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.name = ?", userID, "ROOT").
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// SearchUsers searches for users by company ID
func (r *Repository) SearchUsers(companyID string) ([]struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}, error) {
	var users []struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
	}
	if err := r.db.Model(&model.User{}).
		Select("users.id, users.email").
		Joins("JOIN company_users ON users.id = company_users.user_id").
		Where("company_users.company_id = ?", companyID).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
