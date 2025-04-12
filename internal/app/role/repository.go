package role

import (
	"time"

	"gobizmanager/internal/app/role/permission"
	"gobizmanager/internal/app/user_role"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateAdminRole creates an ADMIN role for a company
func (r *Repository) CreateAdminRole(tx *gorm.DB, companyID int64) (int64, error) {
	role := &Role{
		CompanyID:   companyID,
		Name:        "ADMIN",
		Description: "Company administrator",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := tx.Create(role).Error; err != nil {
		return 0, err
	}
	return role.ID, nil
}

// AssignRoleToUser assigns a role to a user
func (r *Repository) AssignRoleToUser(tx *gorm.DB, userID, roleID int64) error {
	userRole := &user_role.UserRole{
		UserID:    userID,
		RoleID:    roleID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return tx.Create(userRole).Error
}

// GetRoleByID returns a role by ID
func (r *Repository) GetRoleByID(id int64) (*Role, error) {
	var role Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRolesByCompany returns all roles for a company
func (r *Repository) GetRolesByCompany(companyID int64) ([]Role, error) {
	var roles []Role
	if err := r.db.Where("company_id = ?", companyID).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// DeleteRole deletes a role and its associated permissions
func (r *Repository) DeleteRole(tx *gorm.DB, roleID int64) error {
	// First delete user roles
	if err := tx.Where("role_id = ?", roleID).Delete(&user_role.UserRole{}).Error; err != nil {
		return err
	}

	// Then delete permissions
	if err := tx.Where("role_id = ?", roleID).Delete(&permission.Permission{}).Error; err != nil {
		return err
	}

	// Finally delete the role
	return tx.Delete(&Role{}, roleID).Error
}
