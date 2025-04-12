package permission

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreatePermission creates a new permission
func (r *Repository) CreatePermission(name, description string) (int64, error) {
	permission := &Permission{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := r.db.Create(permission).Error; err != nil {
		return 0, fmt.Errorf("failed to create permission: %w", err)
	}

	return permission.ID, nil
}

// AddModuleActionToPermission adds a module action to a permission
func (r *Repository) AddModuleActionToPermission(permissionID, moduleActionID int64) error {
	permissionModuleAction := &PermissionModuleAction{
		PermissionID:   permissionID,
		ModuleActionID: moduleActionID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := r.db.Create(permissionModuleAction).Error; err != nil {
		return fmt.Errorf("failed to add module action to permission: %w", err)
	}

	return nil
}

// AssignPermissionToRole assigns a permission to a role
func (r *Repository) AssignPermissionToRole(roleID, permissionID int64) error {
	rolePermission := &RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := r.db.Create(rolePermission).Error; err != nil {
		return fmt.Errorf("failed to assign permission to role: %w", err)
	}

	return nil
}

// GetRolePermissions returns all permissions for a role
func (r *Repository) GetRolePermissions(roleID int64) ([]Permission, error) {
	var permissions []Permission
	if err := r.db.
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	return permissions, nil
}

// GetPermissionModuleActions returns all module actions for a permission
func (r *Repository) GetPermissionModuleActions(permissionID int64) ([]int64, error) {
	var moduleActionIDs []int64
	if err := r.db.Model(&PermissionModuleAction{}).
		Where("permission_id = ?", permissionID).
		Pluck("module_action_id", &moduleActionIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to get permission module actions: %w", err)
	}

	return moduleActionIDs, nil
}

// GrantAllPermissions grants all module actions to a role
func (r *Repository) GrantAllPermissions(tx *gorm.DB, roleID int64) error {
	var moduleActions []ModuleAction
	if err := tx.Find(&moduleActions).Error; err != nil {
		return err
	}

	for _, ma := range moduleActions {
		// Create permission
		permission := &Permission{
			Name:        ma.Name,
			Description: ma.Description,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := tx.Create(permission).Error; err != nil {
			return err
		}

		// Create permission module action
		permissionModuleAction := &PermissionModuleAction{
			PermissionID:   permission.ID,
			ModuleActionID: ma.ID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := tx.Create(permissionModuleAction).Error; err != nil {
			return err
		}

		// Create role permission
		rolePermission := &RolePermission{
			RoleID:       roleID,
			PermissionID: permission.ID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := tx.Create(rolePermission).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetPermissionsByRoleID returns all permissions for a role
func (r *Repository) GetPermissionsByRoleID(roleID int64) ([]Permission, error) {
	var permissions []Permission
	if err := r.db.
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

// DeletePermissionsByRoleID deletes all permissions for a role
func (r *Repository) DeletePermissionsByRoleID(tx *gorm.DB, roleID int64) error {
	return tx.Where("role_id = ?", roleID).Delete(&Permission{}).Error
}
