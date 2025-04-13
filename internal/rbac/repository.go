package rbac

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"

	model "gobizmanager/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CompanyUser operations
func (r *Repository) CreateCompanyUser(companyID, userID int64, isMain bool) (int64, error) {
	now := time.Now()
	companyUser := &CompanyUser{
		CompanyID: companyID,
		UserID:    userID,
		IsMain:    isMain,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := r.db.Create(companyUser).Error; err != nil {
		return 0, err
	}
	return companyUser.ID, nil
}

func (r *Repository) CreateCompanyUserWithTx(tx *gorm.DB, companyID, userID int64, isMain bool) (int64, error) {
	now := time.Now()
	companyUser := &CompanyUser{
		CompanyID: companyID,
		UserID:    userID,
		IsMain:    isMain,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := tx.Create(companyUser).Error; err != nil {
		return 0, err
	}
	return companyUser.ID, nil
}

func (r *Repository) GetCompanyUserByID(id int64) (*CompanyUser, error) {
	var cu CompanyUser
	if err := r.db.First(&cu, id).Error; err != nil {
		return nil, err
	}
	return &cu, nil
}

func (r *Repository) GetCompanyUserByCompanyAndUser(companyID, userID int64) (*CompanyUser, error) {
	var cu CompanyUser
	if err := r.db.Where("company_id = ? AND user_id = ?", companyID, userID).First(&cu).Error; err != nil {
		return nil, err
	}
	return &cu, nil
}

// Role operations
func (r *Repository) CreateRole(companyID int64, name, description string) (*model.Role, error) {
	now := time.Now()
	role := &model.Role{
		CompanyID:   companyID,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := r.db.Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *Repository) GetRoleByID(id int64) (*model.Role, error) {
	var role model.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// Permission operations
func (r *Repository) CreatePermission(companyID int64, name, description string, roleID int64) (*model.Permission, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer tx.Rollback()

	// Create permission
	permission := &model.Permission{
		CompanyID:   companyID,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := tx.Create(permission).Error; err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	// Associate permission with role
	rolePermission := &RolePermission{
		RoleID:       roleID,
		PermissionID: permission.ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := tx.Create(rolePermission).Error; err != nil {
		return nil, fmt.Errorf("failed to associate permission with role: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return permission, nil
}

func (r *Repository) GetPermissionsByRoleID(roleID int64) ([]model.Permission, error) {
	var permissions []model.Permission
	if err := r.db.
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *Repository) AssignRole(userID, roleID int64) (int64, error) {
	var count int64
	if err := r.db.Model(&model.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, fmt.Errorf("role already assigned to user")
	}

	now := time.Now()
	userRole := &model.UserRole{
		UserID:    userID,
		RoleID:    roleID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := r.db.Create(userRole).Error; err != nil {
		return 0, err
	}
	return userRole.ID, nil
}

func (r *Repository) GetUserRoles(companyUserID int64) ([]model.UserRole, error) {
	var userRoles []model.UserRole
	if err := r.db.Where("company_user_id = ?", companyUserID).Find(&userRoles).Error; err != nil {
		return nil, err
	}
	return userRoles, nil
}

// Module operations
func (r *Repository) GetModuleByID(id int64) (*Module, error) {
	var module Module
	if err := r.db.First(&module, id).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

func (r *Repository) GetModuleActionByID(id int64) (*ModuleAction, error) {
	var moduleAction ModuleAction
	if err := r.db.First(&moduleAction, id).Error; err != nil {
		return nil, err
	}
	return &moduleAction, nil
}

func (r *Repository) HasPermission(userID, moduleActionID int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Permission{}).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Joins("JOIN permission_module_actions ON permissions.id = permission_module_actions.permission_id").
		Where("user_roles.user_id = ? AND permission_module_actions.module_action_id = ?", userID, moduleActionID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) GetUserPermissions(userID int64) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.Model(&model.Permission{}).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *Repository) DeleteCompanyUsersWithTx(tx *gorm.DB, companyID int64) error {
	return tx.Where("company_id = ?", companyID).Delete(&CompanyUser{}).Error
}

func (r *Repository) DeleteCompanyRolesWithTx(tx *gorm.DB, companyID int64) error {
	// First delete user roles
	if err := tx.
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("roles.company_id = ?", companyID).
		Delete(&model.UserRole{}).Error; err != nil {
		return err
	}

	// Then delete role permissions
	if err := tx.
		Joins("JOIN roles ON role_permissions.role_id = roles.id").
		Where("roles.company_id = ?", companyID).
		Delete(&RolePermission{}).Error; err != nil {
		return err
	}

	// Finally delete roles
	return tx.Where("company_id = ?", companyID).Delete(&model.Role{}).Error
}

func (r *Repository) GetCompanyUsersByUserID(userID int64) ([]CompanyUser, error) {
	var companyUsers []CompanyUser
	if err := r.db.Where("user_id = ?", userID).Find(&companyUsers).Error; err != nil {
		return nil, err
	}
	return companyUsers, nil
}

func (r *Repository) CreateRootGroup(name string) (int64, error) {
	rootGroup := &RootGroup{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := r.db.Create(rootGroup).Error; err != nil {
		return 0, err
	}
	return rootGroup.ID, nil
}

func (r *Repository) GetRootGroupByID(id int64) (*RootGroup, error) {
	var rootGroup RootGroup
	if err := r.db.First(&rootGroup, id).Error; err != nil {
		return nil, err
	}
	return &rootGroup, nil
}

func (r *Repository) HasCompanyAccess(userID int64, companyID int64) (bool, error) {
	var count int64
	if err := r.db.Model(&CompanyUser{}).
		Where("user_id = ? AND company_id = ?", userID, companyID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) IsRoot(userID int64) (bool, error) {
	var count int64
	if err := r.db.Model(&model.UserRole{}).
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.name = ? AND roles.company_id IS NULL", userID, "ROOT").
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) AssignPermissionToRole(roleID, permissionID string) error {
	roleIDInt, err := strconv.ParseInt(roleID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid role ID: %w", err)
	}

	permissionIDInt, err := strconv.ParseInt(permissionID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid permission ID: %w", err)
	}

	rolePermission := &RolePermission{
		RoleID:       roleIDInt,
		PermissionID: permissionIDInt,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return r.db.Create(rolePermission).Error
}

func (r *Repository) GetRoleWithPermissions(roleID int64) (*model.Role, error) {
	var role model.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return nil, err
	}

	var permissions []model.Permission
	if err := r.db.
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error; err != nil {
		return nil, err
	}

	role.Permissions = permissions
	return &role, nil
}

func (r *Repository) ListPermissions(companyID int64) ([]model.Permission, error) {
	var permissions []model.Permission
	if err := r.db.Where("company_id = ?", companyID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *Repository) ListRolesWithPermissions(companyID int64) ([]model.Role, error) {
	var roles []model.Role
	if err := r.db.Where("company_id = ?", companyID).Find(&roles).Error; err != nil {
		return nil, err
	}

	for i := range roles {
		var permissions []model.Permission
		if err := r.db.
			Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
			Where("role_permissions.role_id = ?", roles[i].ID).
			Find(&permissions).Error; err != nil {
			return nil, err
		}
		roles[i].Permissions = permissions
	}

	return roles, nil
}

func (r *Repository) RemovePermissionFromRole(roleID int64, permissionID int64) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&RolePermission{}).Error
}

func (r *Repository) GetModuleActionID(module, action string) (int64, error) {
	var moduleAction ModuleAction
	if err := r.db.Where("module = ? AND action = ?", module, action).First(&moduleAction).Error; err != nil {
		return 0, err
	}
	return moduleAction.ID, nil
}

func (r *Repository) GetPermissionByID(id int64) (*model.Permission, error) {
	var permission model.Permission
	if err := r.db.First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *Repository) CreatePermissionModuleAction(permissionID, moduleActionID int64) error {
	return r.db.Create(&PermissionModuleAction{
		PermissionID:   permissionID,
		ModuleActionID: moduleActionID,
	}).Error
}

// Updates the permissions for a role
func (r *Repository) UpdateRolePermissions(roleID string, permissionIDs []int64) error {
	roleIDInt, err := strconv.ParseInt(roleID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid role ID: %w", err)
	}

	if err := r.db.Where("role_id = ?", roleIDInt).Delete(&RolePermission{}).Error; err != nil {
		return err
	}

	for _, permissionID := range permissionIDs {
		rolePermission := RolePermission{
			RoleID:       roleIDInt,
			PermissionID: permissionID,
		}
		if err := r.db.Create(&rolePermission).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) UpdatePermissionModuleActions(permissionID int64, moduleActionIDs []int64) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	// Delete existing permission module actions
	if err := tx.Where("permission_id = ?", permissionID).Delete(&PermissionModuleAction{}).Error; err != nil {
		return err
	}

	// Create new permission module actions
	for _, moduleActionID := range moduleActionIDs {
		permissionModuleAction := &PermissionModuleAction{
			PermissionID:   permissionID,
			ModuleActionID: moduleActionID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := tx.Create(permissionModuleAction).Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (r *Repository) GetModuleActions() ([]ModuleAction, error) {
	var moduleActions []ModuleAction
	if err := r.db.Model(&ModuleAction{}).
		Select("module_actions.id, modules.name as module_name, module_actions.name, module_actions.description").
		Joins("JOIN modules ON module_actions.module_id = modules.id").
		Find(&moduleActions).Error; err != nil {
		return nil, err
	}
	return moduleActions, nil
}

func (r *Repository) GetPermissionModuleActions(permissionID int64) ([]ModuleAction, error) {
	var moduleActions []ModuleAction
	if err := r.db.Model(&ModuleAction{}).
		Select("module_actions.id, modules.name as module_name, module_actions.name, module_actions.description").
		Joins("JOIN modules ON module_actions.module_id = modules.id").
		Joins("JOIN permission_module_actions ON module_actions.id = permission_module_actions.module_action_id").
		Where("permission_module_actions.permission_id = ?", permissionID).
		Find(&moduleActions).Error; err != nil {
		return nil, err
	}
	return moduleActions, nil
}

func (r *Repository) GetCompanyUser(userID int64, companyID int64) (*CompanyUser, error) {
	var companyUser CompanyUser
	if err := r.db.Where("user_id = ? AND company_id = ?", userID, companyID).First(&companyUser).Error; err != nil {
		return nil, err
	}
	return &companyUser, nil
}
