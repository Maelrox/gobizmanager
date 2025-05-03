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
	user := &model.User{
		Email:    username,
		Password: hashedPassword,
		Phone:    phone,
	}
	if err := user.EncryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return 0, err
	}

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

func (r *Repository) GetUserByID(id int64) (*model.User, error) {
	user := &model.User{}
	if err := r.db.First(user, id).Error; err != nil {
		return nil, err
	}

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
	if err := user.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) CreateUser(email, password, phone string) (int64, error) {
	hashedPassword, err := encryption.HashPassword(password)
	if err != nil {
		return 0, err
	}

	user := &model.User{
		Email:    email,
		Password: hashedPassword,
		Phone:    phone,
	}

	if err := user.EncryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return 0, err
	}

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

func (r *Repository) UpdateUser(id int64, email, password, phone string) error {
	user := &model.User{}
	if err := r.db.First(user, id).Error; err != nil {
		return err
	}
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
	if err := user.EncryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return err
	}

	user.UpdatedAt = time.Now()
	return r.db.Save(user).Error
}

func (r *Repository) GetRootRoleID(tx *gorm.DB) (int64, error) {
	var roleID int64
	if err := tx.Model(&model.Role{}).Where("name = ?", "ROOT").Select("id").First(&roleID).Error; err != nil {
		return 0, err
	}
	return roleID, nil
}

func (r *Repository) AssignRootRole(tx *gorm.DB, userID, roleID int64) error {
	userRole := &model.UserRole{
		UserID:    userID,
		RoleID:    roleID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return tx.Create(userRole).Error
}

func (r *Repository) RegisterRootUser(username, password string) (int64, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	userID, err := r.CreateUserWithTx(tx, username, password, "")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	roleID, err := r.GetRootRoleID(tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
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
	var users []model.User
	if err := r.db.Model(&model.User{}).
		Joins("JOIN company_users ON users.id = company_users.user_id").
		Where("company_users.company_id = ?", companyID).
		Find(&users).Error; err != nil {
		return nil, err
	}

	var result []struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
	}
	for _, u := range users {
		if err := u.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
			return nil, err
		}
		result = append(result, struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
		}{
			ID:    uint(u.ID),
			Email: u.Email,
		})
	}
	return result, nil
}
