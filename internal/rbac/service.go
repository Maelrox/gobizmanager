package rbac

import (
	"context"
	"errors"
	"strconv"

	model "gobizmanager/internal/models"
	"gobizmanager/pkg/language"
)

type Service struct {
	repo *Repository
	val  *Validator
}

func NewService(repo *Repository, val *Validator) *Service {
	return &Service{
		repo: repo,
		val:  val,
	}
}

func (s *Service) CreatePermission(ctx context.Context, companyID int64, name, description string, roleID int64) (*model.Permission, error) {
	err := s.val.ValidateCompanyRequest(ctx, companyID)
	if err != nil {
		return nil, err
	}

	role, err := s.repo.GetRoleByID(roleID)
	if err != nil {
		return nil, errors.New(language.RoleNotFound)
	}
	if role.CompanyID != companyID {
		return nil, errors.New(language.RoleNotFound)
	}
	return s.repo.CreatePermission(companyID, name, description, roleID)
}

func (s *Service) AssignRole(ctx context.Context, userID int64, roleID int64) error {
	companyUser, err := s.val.ValidateRoleAssignment(ctx, &AssignRoleRequest{
		UserID: userID,
		RoleID: roleID,
	})
	if err != nil {
		return errors.New(language.PermissionDenied)
	}

	_, err = s.repo.AssignRole(companyUser.ID, roleID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	err := s.val.ValidateRoleRequest(ctx, strconv.FormatInt(roleID, 10))
	if err != nil {
		return err
	}

	err = s.repo.UpdateRolePermissions(strconv.FormatInt(roleID, 10), permissionIDs)
	if err != nil {
		return errors.New(language.PermissionCreateFailed)
	}

	return nil
}

func (s *Service) CreatePermissionModuleAction(ctx context.Context, permissionID int64, moduleActionID int64) error {
	err := s.val.ValidatePermissionRequest(ctx, permissionID)
	if err != nil {
		return err
	}

	return s.repo.CreatePermissionModuleAction(permissionID, moduleActionID)
}

func (s *Service) UpdatePermissionModuleActions(ctx context.Context, permissionID int64, moduleActionIDs []int64) error {
	err := s.val.ValidatePermissionRequest(ctx, permissionID)
	if err != nil {
		return err
	}

	err = s.repo.UpdatePermissionModuleActions(permissionID, moduleActionIDs)
	if err != nil {
		return errors.New(language.PermissionCreateFailed)
	}

	return nil
}

func (s *Service) ListRoles(ctx context.Context, companyID int64) ([]model.Role, error) {
	err := s.val.ValidateCompanyRequest(ctx, companyID)
	if err != nil {
		return nil, err
	}

	roles, err := s.repo.ListRolesWithPermissions(companyID)
	if err != nil {
		return nil, errors.New(language.RoleListFailed)
	}

	return roles, nil
}

func (s *Service) ListPermissions(ctx context.Context, companyID int64) ([]model.Permission, error) {
	err := s.val.ValidateCompanyRequest(ctx, companyID)
	if err != nil {
		return nil, err
	}

	permissions, err := s.repo.ListPermissions(companyID)
	if err != nil {
		return nil, errors.New(language.PermissionListFailed)
	}

	return permissions, nil
}

func (s *Service) GetPermissionModuleActions(ctx context.Context, permissionID int64) ([]ModuleAction, error) {
	err := s.val.ValidatePermissionRequest(ctx, permissionID)
	if err != nil {
		return nil, err
	}

	actions, err := s.repo.GetPermissionModuleActions(permissionID)
	if err != nil {
		return nil, errors.New(language.PermissionListFailed)
	}

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

func (s *Service) CheckPermission(ctx context.Context, userID int64, moduleName, actionName string) (bool, error) {
	moduleActionID, err := s.repo.GetModuleActionID(moduleName, actionName)
	if err != nil {
		return false, errors.New(language.PermissionCheckFailed)
	}
	hasPermission, err := s.repo.HasPermission(userID, moduleActionID)
	if err != nil {
		return false, errors.New(language.PermissionCheckFailed)
	}

	return hasPermission, nil
}

func (s *Service) CreateRole(ctx context.Context, companyID int64, name, description string) (*model.Role, error) {
	err := s.val.ValidateCompanyRequest(ctx, companyID)
	if err != nil {
		return nil, err
	}
	role, err := s.repo.CreateRole(companyID, name, description)
	if err != nil {
		return nil, errors.New(language.RoleCreateFailed)
	}

	return role, nil
}

func (s *Service) RemovePermission(ctx context.Context, roleID, permissionID int64) error {
	err := s.val.ValidateRoleRequest(ctx, strconv.FormatInt(roleID, 10))
	if err != nil {
		return err
	}
	err = s.repo.RemovePermissionFromRole(roleID, permissionID)
	if err != nil {
		return errors.New(language.PermissionRemoveFailed)
	}

	return nil
}

func (s *Service) GetRole(ctx context.Context, roleID int64) (*model.Role, error) {
	err := s.val.ValidateRoleRequest(ctx, strconv.FormatInt(roleID, 10))
	if err != nil {
		return nil, err
	}
	role, err := s.repo.GetRoleWithPermissions(roleID)
	if err != nil {
		return nil, errors.New(language.RoleListFailed)
	}

	return role, nil
}

func (s *Service) CheckRootAccess(ctx context.Context, userID int64) (bool, error) {
	return s.repo.IsRoot(userID)
}

func (s *Service) CheckCompanyAccess(ctx context.Context, userID int64, companyID int64) (bool, error) {
	return s.repo.HasCompanyAccess(userID, companyID)
}

func (s *Service) GetModuleActions(ctx context.Context) ([]ModuleAction, error) {
	return s.repo.GetModuleActions()
}
