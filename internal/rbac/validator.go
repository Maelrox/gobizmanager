package rbac

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"gobizmanager/internal/auth"
	model "gobizmanager/internal/models"
	"gobizmanager/pkg/language"

	"github.com/go-chi/chi/v5"
)

func GetLanguage(ctx context.Context) string {
	if lang, ok := ctx.Value("language").(string); ok {
		return lang
	}
	return "en"
}

type Validator struct {
	Repo     *Repository
	MsgStore *language.MessageStore
}

func NewValidator(repo *Repository, msgStore *language.MessageStore) *Validator {
	return &Validator{
		Repo:     repo,
		MsgStore: msgStore,
	}
}

func (v *Validator) ValidatePermissionRequest(ctx context.Context, permissionID int64) error {

	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return errors.New(language.AuthUserNotFound)
	}

	permission, err := v.Repo.GetPermissionByID(permissionID)
	if err != nil {
		return errors.New(language.PermissionNotFound)
	}

	hasAccess, err := v.Repo.HasCompanyAccess(userID, permission.CompanyID)
	if err != nil {
		return errors.New(language.PermissionCheckFailed)
	}
	if !hasAccess {
		return errors.New(language.PermissionDenied)
	}

	return nil
}

func (v *Validator) ValidateCompanyRequest(ctx context.Context, companyID int64) error {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return errors.New(language.AuthUserNotFound)
	}
	hasAccess, err := v.Repo.HasCompanyAccess(userID, companyID)
	if err != nil {
		return errors.New(language.PermissionCheckFailed)
	}
	if !hasAccess {
		return errors.New(language.PermissionDenied)
	}
	return nil
}

func (v *Validator) ValidateAndGetPermissionID(r *http.Request) (int64, error) {
	permissionID := chi.URLParam(r, "permissionID")
	if permissionID == "" {
		return 0, errors.New(language.ValidationInvalidID)
	}

	permissionIDInt, err := strconv.ParseInt(permissionID, 10, 64)
	if err != nil {
		return 0, errors.New(language.ValidationInvalidID)
	}

	return permissionIDInt, nil
}

func (v *Validator) ValidateAndGetCompanyID(r *http.Request) (int64, error) {
	companyID := chi.URLParam(r, "companyID")
	if companyID == "" {
		return 0, errors.New(language.ValidationInvalidID)
	}

	companyIDInt, err := strconv.ParseInt(companyID, 10, 64)
	if err != nil {
		return 0, errors.New(language.ValidationInvalidID)
	}

	return companyIDInt, nil
}

func (v *Validator) ValidatePermissionAccess(userID int64, permissionID int64) (*model.Permission, error) {
	permission, err := v.Repo.GetPermissionByID(permissionID)
	if err != nil {
		return nil, errors.New(language.PermissionNotFound)
	}

	hasAccess, err := v.Repo.HasCompanyAccess(userID, permission.CompanyID)
	if err != nil {
		return nil, errors.New(language.PermissionCheckFailed)
	}
	if !hasAccess {
		return nil, errors.New(language.PermissionDenied)
	}

	return permission, nil
}

func (v *Validator) ValidateCompanyAccess(userID int64, companyID int64) (bool, error) {
	hasAccess, err := v.Repo.HasCompanyAccess(userID, companyID)
	if err != nil {
		return false, errors.New(language.PermissionCheckFailed)
	}
	if !hasAccess {
		return false, errors.New(language.PermissionDenied)
	}
	return hasAccess, nil
}

func (v *Validator) ValidateRootAccess(userID int64) (bool, error) {
	isRoot, err := v.Repo.IsRoot(userID)
	if err != nil {
		return false, errors.New(language.PermissionCheckFailed)
	}
	if !isRoot {
		return false, errors.New(language.AuthUnauthorized)
	}
	return true, nil
}

func (v *Validator) ValidateAuthenticatedUser(r *http.Request) (int64, error) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		return 0, errors.New(language.AuthUnauthorized)
	}
	return userID, nil
}

func (v *Validator) ValidateRoleAssignment(ctx context.Context, req *AssignRoleRequest) (*CompanyUser, error) {

	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, errors.New(language.AuthUserNotFound)
	}

	role, err := v.Repo.GetRoleByID(req.RoleID)
	if err != nil {
		return nil, errors.New(language.RoleNotFound)
	}

	companyUser, err := v.Repo.GetCompanyUser(userID, role.CompanyID)
	if err != nil {
		return nil, errors.New(language.CompanyUserNotFound)
	}

	if role.Name == "ROOT" {
		return nil, errors.New(language.AuthValidationFailed)
	}

	if role.CompanyID != companyUser.CompanyID {
		return nil, errors.New(language.AuthValidationFailed)
	}

	return companyUser, nil
}

func (v *Validator) ValidateRoleRequest(ctx context.Context, roleID string) error {

	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return errors.New(language.AuthUserNotFound)
	}

	id, err := strconv.ParseInt(roleID, 10, 64)
	if err != nil {
		return errors.New(language.ValidationInvalidID)
	}

	role, err := v.Repo.GetRoleByID(id)
	if err != nil {
		return errors.New(language.RoleNotFound)
	}

	hasAccess, err := v.Repo.HasCompanyAccess(userID, role.CompanyID)
	if err != nil {
		return errors.New(language.PermissionCheckFailed)
	}
	if !hasAccess {
		return errors.New(language.PermissionDenied)
	}

	return nil
}
