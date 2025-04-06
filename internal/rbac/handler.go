package rbac

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"gobizmanager/internal/auth"
	"gobizmanager/pkg/context"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/utils"
)

type Handler struct {
	Repo      *Repository
	Validator *validator.Validate
	MsgStore  *language.MessageStore
}

func NewHandler(repo *Repository, msgStore *language.MessageStore) *Handler {
	return &Handler{
		Repo:      repo,
		Validator: validator.New(),
		MsgStore:  msgStore,
	}
}

// Middleware to check if user has permission
func (h *Handler) RequirePermission(moduleName, actionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := context.GetLanguage(r.Context())

			userID, ok := auth.GetUserID(r.Context())
			if !ok {
				utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, "auth.unauthorized"))
				return
			}

			// Get module action ID
			moduleActionID, err := h.Repo.GetModuleActionID(moduleName, actionName)
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, "rbac.permission_check_failed"))
				return
			}

			hasPermission, err := h.Repo.HasPermission(userID, moduleActionID)
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, "rbac.permission_check_failed"))
				return
			}

			if !hasPermission {
				utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, "rbac.insufficient_permissions"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CreatePermission creates a new permission
func (h *Handler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	// Get user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUserNotFound))
		return
	}

	// Get company ID from URL
	companyID := chi.URLParam(r, "companyID")
	if companyID == "" {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationInvalidID))
		return
	}

	// Verify user has access to this company
	hasAccess, err := h.Repo.HasCompanyAccess(userID, companyID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !hasAccess {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	var req CreatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	// Check if the target role is ROOT
	role, err := h.Repo.GetRoleByID(req.RoleID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCreateFailed))
		return
	}

	if role.Name == "ROOT" {
		utils.JSONError(w, http.StatusForbidden, "Cannot add permissions to ROOT role")
		return
	}

	permission, err := h.Repo.CreatePermission(req.Name, req.Description, req.RoleID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCreateFailed))
		return
	}

	utils.JSON(w, http.StatusCreated, permission)
}

// CreateRole creates a new role with permissions
func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	// Get user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUserNotFound))
		return
	}

	// Check if user is ROOT
	isRoot, err := h.Repo.IsRoot(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !isRoot {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	var req CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	role, err := h.Repo.CreateRoleWithPermissions(req.Name, req.Description, req.Permissions)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgRoleCreateFailed))
		return
	}

	utils.JSON(w, http.StatusCreated, role)
}

// AssignPermission assigns a permission to a role
func (h *Handler) AssignPermission(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	// Get user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUserNotFound))
		return
	}

	// Check if user is ROOT
	isRoot, err := h.Repo.IsRoot(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !isRoot {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	var req AssignPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	if err := h.Repo.AssignPermissionToRole(req.RoleID, req.PermissionID); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionAssignFailed))
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"message": h.MsgStore.GetMessage(lang, language.MsgPermissionAssigned),
	})
}

// RemovePermission removes a permission from a role
func (h *Handler) RemovePermission(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	// Get user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUserNotFound))
		return
	}

	// Check if user is ROOT
	isRoot, err := h.Repo.IsRoot(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !isRoot {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	var req RemovePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	if err := h.Repo.RemovePermissionFromRole(req.RoleID, req.PermissionID); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionRemoveFailed))
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"message": h.MsgStore.GetMessage(lang, language.MsgPermissionRemoved),
	})
}

// GetRole retrieves a role with its permissions
func (h *Handler) GetRole(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	// Get user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUserNotFound))
		return
	}

	// Check if user is ROOT
	isRoot, err := h.Repo.IsRoot(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !isRoot {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	roleID := chi.URLParam(r, "id")
	if roleID == "" {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationInvalidID))
		return
	}

	role, err := h.Repo.GetRoleWithPermissions(roleID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgRoleNotFound))
		return
	}

	utils.JSON(w, http.StatusOK, role)
}

// ListRoles retrieves all roles with their permissions
func (h *Handler) ListRoles(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	// Get user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUserNotFound))
		return
	}

	// Check if user is ROOT
	isRoot, err := h.Repo.IsRoot(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !isRoot {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	roles, err := h.Repo.ListRolesWithPermissions()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgRoleListFailed))
		return
	}

	utils.JSON(w, http.StatusOK, roles)
}

// ListPermissions retrieves all permissions for a company
func (h *Handler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	// Get user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUserNotFound))
		return
	}

	// Get company ID from URL
	companyID := chi.URLParam(r, "companyID")
	if companyID == "" {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationInvalidID))
		return
	}

	companyIDInt, err := strconv.ParseInt(companyID, 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationInvalidID))
		return
	}

	// Verify user has access to this company
	hasAccess, err := h.Repo.HasCompanyAccess(userID, companyID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !hasAccess {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	permissions, err := h.Repo.ListPermissions(companyIDInt)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionListFailed))
		return
	}

	utils.JSON(w, http.StatusOK, permissions)
}

// AssignRole assigns a role to a company user
func (h *Handler) AssignRole(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	// Get user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUserNotFound))
		return
	}

	var req AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	// Get the company user to verify it exists and get the company ID
	companyUsers, err := h.Repo.GetCompanyUsersByUserID(req.UserID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if len(companyUsers) == 0 {
		utils.JSONError(w, http.StatusNotFound, h.MsgStore.GetMessage(lang, language.MsgCompanyUserNotFound))
		return
	}

	// For now, we'll use the first company user found
	// TODO: Consider if we need to handle multiple company relationships
	companyUser := companyUsers[0]

	// Verify the requesting user has access to this company
	hasAccess, err := h.Repo.HasCompanyAccess(userID, strconv.FormatInt(companyUser.CompanyID, 10))
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !hasAccess {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	// Get the role to verify it's not ROOT and belongs to the same company
	role, err := h.Repo.GetRoleByID(req.RoleID)
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, h.MsgStore.GetMessage(lang, language.MsgRoleNotFound))
		return
	}

	// Check if role is ROOT
	if role.Name == "ROOT" {
		utils.JSONError(w, http.StatusForbidden, "Cannot assign ROOT role")
		return
	}

	// Verify role belongs to the same company
	if role.CompanyID != companyUser.CompanyID {
		utils.JSONError(w, http.StatusForbidden, "Role does not belong to the company")
		return
	}

	_, err = h.Repo.AssignRole(companyUser.ID, req.RoleID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgRoleAssignFailed))
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"message": h.MsgStore.GetMessage(lang, language.MsgRoleAssigned),
	})
}
