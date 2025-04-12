package rbac

import (
	"context"
	"net/http"
	"strconv"

	"gobizmanager/internal/app/auth"
	pkgctx "gobizmanager/pkg/context"
	"gobizmanager/pkg/errors"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/utils"

	"github.com/go-chi/chi/v5"
)

// GetLanguage returns the language from the context
func GetLanguage(ctx context.Context) string {
	if lang, ok := ctx.Value("language").(string); ok {
		return lang
	}
	return "en" // default language
}

// Validator handles all validation logic
type Validator struct {
	Repo     *Repository
	MsgStore *language.MessageStore
}

// ValidationResult holds the result of a validation
type ValidationResult struct {
	UserID       int64
	PermissionID int64
	CompanyID    int64
	Permission   *Permission
	Role         *Role
	HasAccess    bool
	CompanyUser  CompanyUser
}

func NewValidator(repo *Repository, msgStore *language.MessageStore) *Validator {
	return &Validator{
		Repo:     repo,
		MsgStore: msgStore,
	}
}

// HandleValidationError is a middleware to handle validation errors
func (v *Validator) HandleValidationError(w http.ResponseWriter, r *http.Request, next func()) {
	lang := pkgctx.GetLanguage(r.Context())

	defer func() {
		if err := recover(); err != nil {
			if appErr, ok := err.(*errors.Error); ok {
				utils.JSONError(w, appErr.GetHTTPStatus(), v.MsgStore.GetMessage(lang, appErr.Message))
			} else {
				panic(err)
			}
		}
	}()

	next()
}

// ValidatePermissionRequest validates all permission-related requests
func (v *Validator) ValidatePermissionRequest(ctx context.Context, permissionID string) *ValidationResult {
	result := &ValidationResult{}

	// Get user ID from context
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		panic(NewError(ErrorTypeAuthentication, ErrorCodeUserNotFound, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgAuthUserNotFound)))
	}
	result.UserID = userID

	// Parse permission ID
	id, err := strconv.ParseInt(permissionID, 10, 64)
	if err != nil {
		panic(NewError(ErrorTypeValidation, ErrorCodeInvalidID, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgValidationInvalidID)))
	}
	result.PermissionID = id

	// Get permission
	permission, err := v.Repo.GetPermissionByID(id)
	if err != nil {
		panic(NewError(ErrorTypeNotFound, ErrorCodePermissionNotFound, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgPermissionNotFound)))
	}
	result.Permission = permission

	// Verify user has access to this company
	hasAccess, err := v.Repo.HasCompanyAccess(userID, strconv.FormatInt(permission.CompanyID, 10))
	if err != nil {
		panic(NewError(ErrorTypeInternal, ErrorCodePermissionCheckFailed, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgPermissionCheckFailed)))
	}
	if !hasAccess {
		panic(NewError(ErrorTypeAuthorization, ErrorCodePermissionDenied, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgPermissionDenied)))
	}

	return result
}

// ValidateCompanyRequest validates all company-related requests
func (v *Validator) ValidateCompanyRequest(ctx context.Context, companyID string) *ValidationResult {
	result := &ValidationResult{}

	// Get user ID from context
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		panic(NewError(ErrorTypeAuthentication, ErrorCodeUserNotFound, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgAuthUserNotFound)))
	}
	result.UserID = userID

	// Parse company ID
	id, err := strconv.ParseInt(companyID, 10, 64)
	if err != nil {
		panic(NewError(ErrorTypeValidation, ErrorCodeInvalidID, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgValidationInvalidID)))
	}
	result.CompanyID = id

	// Verify user has access to this company
	hasAccess, err := v.Repo.HasCompanyAccess(userID, companyID)
	if err != nil {
		panic(NewError(ErrorTypeInternal, ErrorCodePermissionCheckFailed, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgPermissionCheckFailed)))
	}
	if !hasAccess {
		panic(NewError(ErrorTypeAuthorization, ErrorCodePermissionDenied, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgPermissionDenied)))
	}

	return result
}

// ValidateRootRequest validates all root-level requests
func (v *Validator) ValidateRootRequest(r *http.Request) ValidationResult {
	result := ValidationResult{
		UserID: v.ValidateAuthenticatedUser(r),
	}
	result.HasAccess = v.ValidateRootAccess(result.UserID)
	return result
}

// ValidateAndGetPermissionID validates and gets permission ID from URL
func (v *Validator) ValidateAndGetPermissionID(r *http.Request) int64 {
	permissionID := chi.URLParam(r, "permissionID")
	if permissionID == "" {
		panic(NewError(ErrorTypeValidation, ErrorCodeInvalidID, v.MsgStore.GetMessage(GetLanguage(r.Context()), language.MsgValidationInvalidID)))
	}

	permissionIDInt, err := strconv.ParseInt(permissionID, 10, 64)
	if err != nil {
		panic(NewError(ErrorTypeValidation, ErrorCodeInvalidID, v.MsgStore.GetMessage(GetLanguage(r.Context()), language.MsgValidationInvalidID)))
	}

	return permissionIDInt
}

// ValidateAndGetCompanyID validates and gets company ID from URL
func (v *Validator) ValidateAndGetCompanyID(r *http.Request) int64 {
	companyID := chi.URLParam(r, "companyID")
	if companyID == "" {
		panic(NewError(ErrorTypeValidation, ErrorCodeInvalidID, v.MsgStore.GetMessage(GetLanguage(r.Context()), language.MsgValidationInvalidID)))
	}

	companyIDInt, err := strconv.ParseInt(companyID, 10, 64)
	if err != nil {
		panic(NewError(ErrorTypeValidation, ErrorCodeInvalidID, v.MsgStore.GetMessage(GetLanguage(r.Context()), language.MsgValidationInvalidID)))
	}

	return companyIDInt
}

// ValidatePermissionAccess validates if user has access to a permission
func (v *Validator) ValidatePermissionAccess(userID int64, permissionID int64) *Permission {
	permission, err := v.Repo.GetPermissionByID(permissionID)
	if err != nil {
		panic(NewError(ErrorTypeNotFound, ErrorCodePermissionNotFound, v.MsgStore.GetMessage("en", language.MsgPermissionNotFound)))
	}

	hasAccess, err := v.Repo.HasCompanyAccess(userID, strconv.FormatInt(permission.CompanyID, 10))
	if err != nil {
		panic(NewError(ErrorTypeInternal, ErrorCodePermissionCheckFailed, v.MsgStore.GetMessage("en", language.MsgPermissionCheckFailed)))
	}
	if !hasAccess {
		panic(NewError(ErrorTypeAuthorization, ErrorCodePermissionDenied, v.MsgStore.GetMessage("en", language.MsgPermissionDenied)))
	}

	return permission
}

// ValidateCompanyAccess validates if user has access to a company
func (v *Validator) ValidateCompanyAccess(userID int64, companyID string) bool {
	hasAccess, err := v.Repo.HasCompanyAccess(userID, companyID)
	if err != nil {
		panic(NewError(ErrorTypeInternal, ErrorCodePermissionCheckFailed, v.MsgStore.GetMessage("en", language.MsgPermissionCheckFailed)))
	}
	if !hasAccess {
		panic(NewError(ErrorTypeAuthorization, ErrorCodePermissionDenied, v.MsgStore.GetMessage("en", language.MsgPermissionDenied)))
	}
	return hasAccess
}

// ValidateRootAccess validates if user has root access
func (v *Validator) ValidateRootAccess(userID int64) bool {
	isRoot, err := v.Repo.IsRoot(userID)
	if err != nil {
		panic(NewError(ErrorTypeInternal, ErrorCodePermissionCheckFailed, v.MsgStore.GetMessage("en", language.MsgPermissionCheckFailed)))
	}
	if !isRoot {
		panic(NewError(ErrorTypeAuthorization, ErrorCodeRootAccessDenied, v.MsgStore.GetMessage("en", MsgRootAccessDenied)))
	}
	return isRoot
}

// ValidateAuthenticatedUser validates and gets the authenticated user
func (v *Validator) ValidateAuthenticatedUser(r *http.Request) int64 {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		panic(NewError(ErrorTypeAuthentication, ErrorCodeUnauthorized, v.MsgStore.GetMessage(GetLanguage(r.Context()), language.MsgAuthUnauthorized)))
	}
	return userID
}

// ValidateRoleAssignment validates all role assignment requests
func (v *Validator) ValidateRoleAssignment(ctx context.Context, req *AssignRoleRequest) *ValidationResult {
	result := &ValidationResult{}

	// Get user ID from context
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		panic(NewError(ErrorTypeAuthentication, ErrorCodeUserNotFound, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgAuthUserNotFound)))
	}
	result.UserID = userID

	// Get role first
	role, err := v.Repo.GetRoleByID(req.RoleID)
	if err != nil {
		panic(NewError(ErrorTypeNotFound, ErrorCodeRoleNotFound, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgRoleNotFound)))
	}
	result.Role = role

	// Get company user using role's company ID
	companyUser, err := v.Repo.GetCompanyUser(userID, strconv.FormatInt(role.CompanyID, 10))
	if err != nil {
		panic(NewError(ErrorTypeNotFound, ErrorCodeCompanyUserNotFound, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgCompanyUserNotFound)))
	}
	result.CompanyUser = *companyUser

	// Verify role is not ROOT
	if role.Name == "ROOT" {
		panic(NewError(ErrorTypeAuthorization, ErrorCodeRootRoleAssignment, v.MsgStore.GetMessage(GetLanguage(ctx), MsgRootRoleAssignment)))
	}

	// Verify role belongs to the same company
	if role.CompanyID != companyUser.CompanyID {
		panic(NewError(ErrorTypeAuthorization, ErrorCodeRoleCompanyMismatch, v.MsgStore.GetMessage(GetLanguage(ctx), MsgRoleCompanyMismatch)))
	}

	return result
}

// ValidateRoleRequest validates a role request
func (v *Validator) ValidateRoleRequest(ctx context.Context, roleID string) *ValidationResult {
	result := &ValidationResult{}

	// Get user ID from context
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		panic(NewError(ErrorTypeAuthentication, ErrorCodeUserNotFound, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgAuthUserNotFound)))
	}
	result.UserID = userID

	// Parse role ID
	id, err := strconv.ParseInt(roleID, 10, 64)
	if err != nil {
		panic(NewError(ErrorTypeValidation, ErrorCodeInvalidID, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgValidationInvalidID)))
	}

	// Get role
	role, err := v.Repo.GetRoleByID(id)
	if err != nil {
		panic(NewError(ErrorTypeNotFound, ErrorCodeRoleNotFound, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgRoleNotFound)))
	}
	result.Role = role

	// Verify user has access to this company
	hasAccess, err := v.Repo.HasCompanyAccess(userID, strconv.FormatInt(role.CompanyID, 10))
	if err != nil {
		panic(NewError(ErrorTypeInternal, ErrorCodePermissionCheckFailed, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgPermissionCheckFailed)))
	}
	if !hasAccess {
		panic(NewError(ErrorTypeAuthorization, ErrorCodePermissionDenied, v.MsgStore.GetMessage(GetLanguage(ctx), language.MsgPermissionDenied)))
	}

	return result
}
