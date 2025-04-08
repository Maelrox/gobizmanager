package rbac

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	pkgctx "gobizmanager/pkg/context"
	"gobizmanager/pkg/errors"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/utils"
)

// RoleHandler handles all role-related HTTP requests
type RoleHandler struct {
	*BaseHandler
}

func NewRoleHandler(repo *Repository, msgStore *language.MessageStore) *RoleHandler {
	return &RoleHandler{
		BaseHandler: NewBaseHandler(repo, msgStore),
	}
}

// CreateRole creates a new role
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	h.HandleValidationError(w, r, func() {
		var req CreateRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		role, err := h.Service.CreateRole(r.Context(), req.CompanyID, req.Name, req.Description)
		if err != nil {
			panic(err)
		}

		utils.JSON(w, http.StatusCreated, role)
	})
}

// GetRole returns a role by ID
func (h *RoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	h.HandleValidationError(w, r, func() {
		roleID := chi.URLParam(r, "id")
		if roleID == "" {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		roleIDInt, err := strconv.ParseInt(roleID, 10, 64)
		if err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		role, err := h.Service.GetRole(r.Context(), roleIDInt)
		if err != nil {
			panic(err)
		}

		utils.JSON(w, http.StatusOK, role)
	})
}

// ListRoles returns all roles for a company
func (h *RoleHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	h.HandleValidationError(w, r, func() {
		companyID := chi.URLParam(r, "companyID")
		companyIDInt, err := strconv.ParseInt(companyID, 10, 64)
		if err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		roles, err := h.Service.ListRoles(r.Context(), companyIDInt)
		if err != nil {
			panic(err)
		}

		utils.JSON(w, http.StatusOK, roles)
	})
}

// AssignRole assigns a role to a user
func (h *RoleHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	h.HandleValidationError(w, r, func() {
		var req AssignRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		if err := h.Service.AssignRole(r.Context(), req.UserID, req.RoleID); err != nil {
			panic(err)
		}

		utils.JSON(w, http.StatusOK, map[string]interface{}{
			"message": h.MsgStore.GetMessage(pkgctx.GetLanguage(r.Context()), language.MsgRoleAssigned),
		})
	})
}

// UpdateRolePermissions updates permissions for a role
func (h *RoleHandler) UpdateRolePermissions(w http.ResponseWriter, r *http.Request) {
	h.HandleValidationError(w, r, func() {
		var req UpdateRolePermissionsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		roleID, err := strconv.ParseInt(req.RoleID, 10, 64)
		if err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		if err := h.Service.UpdateRolePermissions(r.Context(), roleID, req.PermissionIDs); err != nil {
			panic(err)
		}

		utils.JSON(w, http.StatusOK, map[string]interface{}{
			"message": h.MsgStore.GetMessage(pkgctx.GetLanguage(r.Context()), language.MsgRoleAssigned),
		})
	})
}
