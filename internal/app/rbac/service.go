package rbac

import (
	"context"
	"strconv"

	"gobizmanager/pkg/errors"
	"gobizmanager/pkg/language"
)

// Service handles RBAC business logic
type Service struct {
	repo *Repository
	val  *Validator
}

// NewService creates a new RBAC service
func NewService(repo *Repository, val *Validator) *Service {
	return &Service{
		repo: repo,
		val:  val,
	}
}

func (s *Service) CreatePermission(ctx context.Context, companyID int64, name, description string, roleID int64) (*Permission, error) {
	// Validate company access
	result := s.val.ValidateCompanyRequest(ctx, strconv.FormatInt(companyID, 10))
	if result.CompanyID != companyID {
		return nil, s.newError(ctx, errors.ErrorTypeAuthorization, errors.ErrorCodePermissionDenied, language.MsgPermissionDenied)
	}

	// Check if role exists and belongs to the company
	role, err := s.repo.GetRoleByID(roleID)
	if err != nil {
		return nil, s.newError(ctx, errors.ErrorTypeNotFound, errors.ErrorCodeRoleNotFound, language.MsgRoleNotFound)
	}
	if role.CompanyID != companyID {
		return nil, s.newError(ctx, errors.ErrorTypeAuthorization, errors.ErrorCodeRoleCompanyMismatch, MsgRoleCompanyMismatch)
	}
	return s.repo.CreatePermission(companyID, name, description, roleID)
}

// AssignRole assigns a role to a user
func (s *Service) AssignRole(ctx context.Context, userID int64, roleID int64) error {
	result := s.val.ValidateRoleAssignment(ctx, &AssignRoleRequest{
		UserID: userID,
		RoleID: roleID,
	})

	// Assign role
	_, err := s.repo.AssignRole(result.CompanyUser.ID, roleID)
	return err
}

// UpdateRolePermissions updates permissions for a role
func (s *Service) UpdateRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	// Validate role request
	result := s.val.ValidateRoleRequest(ctx, strconv.FormatInt(roleID, 10))
	if result.Role.ID != strconv.FormatInt(roleID, 10) {
		return s.newError(ctx, errors.ErrorTypeAuthorization, errors.ErrorCodePermissionDenied, language.MsgPermissionDenied)
	}

	// Update permissions
	return s.repo.UpdateRolePermissions(strconv.FormatInt(roleID, 10), permissionIDs)
}

// CreatePermissionModuleAction associates a module action with a permission
func (s *Service) CreatePermissionModuleAction(ctx context.Context, permissionID int64, moduleActionID int64) error {
	// Validate permission request
	result := s.val.ValidatePermissionRequest(ctx, strconv.FormatInt(permissionID, 10))
	if result.PermissionID != permissionID {
		return s.newError(ctx, errors.ErrorTypeAuthorization, errors.ErrorCodePermissionDenied, language.MsgPermissionDenied)
	}

	// Create association
	return s.repo.CreatePermissionModuleAction(permissionID, moduleActionID)
}

// UpdatePermissionModuleActions updates module actions for a permission
func (s *Service) UpdatePermissionModuleActions(ctx context.Context, permissionID int64, moduleActionIDs []int64) error {
	// Validate permission request
	result := s.val.ValidatePermissionRequest(ctx, strconv.FormatInt(permissionID, 10))
	if result.PermissionID != permissionID {
		return s.newError(ctx, errors.ErrorTypeAuthorization, errors.ErrorCodePermissionDenied, language.MsgPermissionDenied)
	}

	// Update actions
	return s.repo.UpdatePermissionModuleActions(strconv.FormatInt(permissionID, 10), moduleActionIDs)
}

// ListRoles returns all roles for a company
func (s *Service) ListRoles(ctx context.Context, companyID int64) ([]Role, error) {
	// Validate company request
	result := s.val.ValidateCompanyRequest(ctx, strconv.FormatInt(companyID, 10))
	if result.CompanyID != companyID {
		return nil, s.newError(ctx, errors.ErrorTypeAuthorization, errors.ErrorCodePermissionDenied, language.MsgPermissionDenied)
	}

	// List roles
	return s.repo.ListRolesWithPermissions(companyID)
}

// ListPermissions returns all permissions for a company
func (s *Service) ListPermissions(ctx context.Context, companyID int64) ([]Permission, error) {
	// Validate company request
	result := s.val.ValidateCompanyRequest(ctx, strconv.FormatInt(companyID, 10))
	if result.CompanyID != companyID {
		return nil, s.newError(ctx, errors.ErrorTypeAuthorization, errors.ErrorCodePermissionDenied, language.MsgPermissionDenied)
	}

	// List permissions
	return s.repo.ListPermissions(companyID)
}

// GetPermissionModuleActions returns all available actions for a permission module
func (s *Service) GetPermissionModuleActions(ctx context.Context, permissionID int64) ([]ModuleAction, error) {
	// Validate permission request
	result := s.val.ValidatePermissionRequest(ctx, strconv.FormatInt(permissionID, 10))
	if result.PermissionID != permissionID {
		return nil, s.newError(ctx, errors.ErrorTypeAuthorization, errors.ErrorCodePermissionDenied, language.MsgPermissionDenied)
	}

	// Get actions
	actions, err := s.repo.GetPermissionModuleActions(permissionID)
	if err != nil {
		return nil, err
	}

	// Convert to ModuleAction slice
	moduleActions := make([]ModuleAction, len(actions))
	for i, action := range actions {
		moduleActions[i] = ModuleAction{
			ID:          action.ID,
			Name:        action.Name,
			Description: action.Description,
		}
	}

	return moduleActions, nil
}

// CheckPermission checks if a user has a specific permission
func (s *Service) CheckPermission(ctx context.Context, userID int64, moduleName, actionName string) (bool, error) {
	// Get module action ID
	moduleActionID, err := s.repo.GetModuleActionID(moduleName, actionName)
	if err != nil {
		return false, s.newError(ctx, errors.ErrorTypeInternal, errors.ErrorCodePermissionCheckFailed, language.MsgPermissionCheckFailed)
	}

	// Check permission
	return s.repo.HasPermission(userID, moduleActionID)
}

// CreateRole creates a new role
func (s *Service) CreateRole(ctx context.Context, companyID int64, name, description string) (*Role, error) {
	// Validate company access
	result := s.val.ValidateCompanyRequest(ctx, strconv.FormatInt(companyID, 10))
	if result.CompanyID != companyID {
		return nil, s.newError(ctx, errors.ErrorTypeAuthorization, errors.ErrorCodePermissionDenied, language.MsgPermissionDenied)
	}

	// Create role
	return s.repo.CreateRole(companyID, name, description)
}

// RemovePermission removes a permission from a role
func (s *Service) RemovePermission(ctx context.Context, roleID, permissionID int64) error {
	// Validate role request
	result := s.val.ValidateRoleRequest(ctx, strconv.FormatInt(roleID, 10))
	if result.Role.ID != strconv.FormatInt(roleID, 10) {
		return s.newError(ctx, errors.ErrorTypeAuthorization, errors.ErrorCodePermissionDenied, language.MsgPermissionDenied)
	}

	// Remove permission
	return s.repo.RemovePermissionFromRole(strconv.FormatInt(roleID, 10), strconv.FormatInt(permissionID, 10))
}

// GetRole returns a role by ID
func (s *Service) GetRole(ctx context.Context, roleID int64) (*Role, error) {
	// Validate role request
	result := s.val.ValidateRoleRequest(ctx, strconv.FormatInt(roleID, 10))
	if result.Role.ID != strconv.FormatInt(roleID, 10) {
		return nil, s.newError(ctx, errors.ErrorTypeAuthorization, errors.ErrorCodePermissionDenied, language.MsgPermissionDenied)
	}

	// Get role
	return s.repo.GetRoleWithPermissions(strconv.FormatInt(roleID, 10))
}

// CheckRootAccess checks if a user has root access
func (s *Service) CheckRootAccess(ctx context.Context, userID int64) (bool, error) {
	return s.repo.IsRoot(userID)
}

// CheckCompanyAccess checks if a user has access to a company
func (s *Service) CheckCompanyAccess(ctx context.Context, userID int64, companyID string) (bool, error) {
	return s.repo.HasCompanyAccess(userID, companyID)
}

// GetModuleActions returns all module actions
func (s *Service) GetModuleActions(ctx context.Context) ([]struct {
	ID          int64  `json:"id"`
	ModuleName  string `json:"module_name"`
	Name        string `json:"name"`
	Description string `json:"description"`
}, error) {
	return s.repo.GetModuleActions()
}
