package company

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"gobizmanager/internal/app/rbac"
	"gobizmanager/pkg/language"
	"gobizmanager/platform/config"

	"gorm.io/gorm"
)

type Repository struct {
	db       *gorm.DB
	cfg      *config.Config
	RBACRepo *rbac.Repository
}

func NewRepository(db *gorm.DB, cfg *config.Config, rbacRepo *rbac.Repository) *Repository {
	return &Repository{
		db:       db,
		cfg:      cfg,
		RBACRepo: rbacRepo,
	}
}

func (r *Repository) CreateCompany(req *CreateCompanyRequest, userID int64) (*Company, error) {
	company := &Company{
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Address:    req.Address,
		Logo:       sql.NullString{String: req.Logo, Valid: req.Logo != ""},
		Identifier: req.Identifier,
	}
	if err := company.EncryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return nil, fmt.Errorf("failed to encrypt company fields: %w", err)
	}

	// Start transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create company
	if err := tx.Create(company).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create company: %w", err)
	}

	// Create company-user relationship
	companyUser := &rbac.CompanyUser{
		CompanyID: company.ID,
		UserID:    userID,
		IsMain:    true,
	}
	if err := tx.Create(companyUser).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create company-user relationship: %w", err)
	}

	// Create ADMIN role for the company
	adminRole := &rbac.Role{
		CompanyID:   company.ID,
		Name:        "ADMIN",
		Description: "Company administrator",
	}
	if err := tx.Create(adminRole).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create ADMIN role: %w", err)
	}

	// Get all generic permissions
	var permissions []rbac.Permission
	if err := tx.Where("company_id IS NULL").Find(&permissions).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	// Assign all permissions to ADMIN role
	for _, permission := range permissions {
		rolePermission := &rbac.Permission{
			CompanyID:   company.ID,
			Name:        permission.Name,
			Description: permission.Description,
		}
		if err := tx.Create(rolePermission).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to assign permission to ADMIN role: %w", err)
		}
	}

	// Convert role ID from string to int64
	roleID, err := strconv.ParseInt(adminRole.ID, 10, 64)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to parse role ID: %w", err)
	}

	// Assign ADMIN role to user
	userRole := &rbac.UserRole{
		CompanyUserID: companyUser.ID,
		RoleID:        roleID,
	}
	if err := tx.Create(userRole).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to assign ADMIN role to user: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	if err := company.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}
	return company, nil
}

func (r *Repository) GetCompany(id int64) (*Company, error) {
	var company Company
	if err := r.db.First(&company, id).Error; err != nil {
		return nil, err
	}

	// Decrypt sensitive fields
	if err := company.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return nil, err
	}

	return &company, nil
}

func (r *Repository) UpdateCompany(id int64, req *UpdateCompanyRequest) (*Company, error) {
	// Create a temporary company to encrypt fields
	company := &Company{
		Name:       req.Name,
		Phone:      req.Phone,
		Email:      req.Email,
		Identifier: req.Identifier,
	}

	// Encrypt sensitive fields
	if err := company.EncryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return nil, err
	}

	if err := r.db.Model(&Company{}).Where("id = ?", id).Updates(company).Error; err != nil {
		return nil, err
	}

	if err := company.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}
	return r.GetCompany(id)
}

func (r *Repository) UpdateCompanyLogo(id string, logo string) error {
	return r.db.Model(&Company{}).Where("id = ?", id).Update("logo", logo).Error
}

func (r *Repository) DeleteCompanyWithTx(tx *gorm.DB, companyID int64) error {
	// Delete company-user relationships
	if err := r.RBACRepo.DeleteCompanyUsersWithTx(tx, companyID); err != nil {
		return err
	}

	// Delete company roles and permissions
	if err := r.RBACRepo.DeleteCompanyRolesWithTx(tx, companyID); err != nil {
		return err
	}

	// Delete the company
	return tx.Delete(&Company{}, companyID).Error
}

func (r *Repository) ListCompanies(userID int64) ([]*Company, error) {
	var companies []*Company
	err := r.db.Joins("JOIN company_users ON companies.id = company_users.company_id").
		Where("company_users.user_id = ?", userID).
		Find(&companies).Error
	if err != nil {
		return nil, err
	}

	// Decrypt sensitive fields for each company
	for _, company := range companies {
		if err := company.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
			return nil, err
		}
	}

	return companies, nil
}

// ListCompaniesForUser returns all companies a user has access to
func (r *Repository) ListCompaniesForUser(userID int64) ([]Company, error) {
	var companies []Company
	err := r.db.Joins("JOIN company_users ON companies.id = company_users.company_id").
		Where("company_users.user_id = ?", userID).
		Find(&companies).Error
	if err != nil {
		return nil, err
	}

	// Decrypt sensitive fields for each company
	for i := range companies {
		if err := companies[i].DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
			return nil, err
		}
	}

	return companies, nil
}

// CompanyExistsForUser checks if a company with the given name already exists for the user
func (r *Repository) CompanyExistsForUser(userID int64, name string) (bool, error) {
	var count int64
	err := r.db.Model(&Company{}).
		Joins("JOIN company_users ON companies.id = company_users.company_id").
		Where("company_users.user_id = ? AND companies.name = ?", userID, name).
		Count(&count).Error
	if err != nil {
		return false, errors.New(language.CompanyListFailed)
	}
	return count > 0, nil
}

// CompanyExistsForUser checks if a company with the given name already exists for the user
func (r *Repository) CompanyExistsForUserByID(userID int64, companyID int64) (bool, error) {
	var count int64
	err := r.db.Model(&Company{}).
		Joins("JOIN company_users ON companies.id = company_users.company_id").
		Where("company_users.user_id = ? AND company_users.companies.id = ?", userID, companyID).
		Count(&count).Error
	if err != nil {
		return false, errors.New(language.CompanyListFailed)
	}
	return count > 0, nil
}

// AddUserToCompany adds a user to a company
func (r *Repository) AddUserToCompany(tx *gorm.DB, companyID, userID int64) error {
	companyUser := &rbac.CompanyUser{
		CompanyID: companyID,
		UserID:    userID,
	}
	return tx.Create(companyUser).Error
}

func (r *Repository) DeleteCompany(companyID int64) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := r.DeleteCompanyWithTx(tx, companyID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
