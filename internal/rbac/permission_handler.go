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

// PermissionHandler handles all permission-related HTTP requests
type PermissionHandler struct {
	*BaseHandler
}

func NewPermissionHandler(repo *Repository, msgStore *language.MessageStore) *PermissionHandler {
	return &PermissionHandler{
		BaseHandler: NewBaseHandler(repo, msgStore),
	}
}

// CreatePermission creates a new permission
func (h *PermissionHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	h.HandleValidationError(w, r, func() {
		var req CreatePermissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		permission, err := h.Service.CreatePermission(r.Context(), req.CompanyID, req.Name, req.Description, req.RoleID)
		if err != nil {
			panic(err)
		}

		utils.JSON(w, http.StatusCreated, permission)
	})
}

// ListPermissions returns all permissions for a company
func (h *PermissionHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	h.HandleValidationError(w, r, func() {
		companyID := chi.URLParam(r, "companyID")
		companyIDInt, err := strconv.ParseInt(companyID, 10, 64)
		if err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		permissions, err := h.Service.ListPermissions(r.Context(), companyIDInt)
		if err != nil {
			panic(err)
		}

		utils.JSON(w, http.StatusOK, permissions)
	})
}

// RemovePermission removes a permission from a role
func (h *PermissionHandler) RemovePermission(w http.ResponseWriter, r *http.Request) {
	h.HandleValidationError(w, r, func() {
		var req RemovePermissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		roleID, err := strconv.ParseInt(req.RoleID, 10, 64)
		if err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		permissionID, err := strconv.ParseInt(req.PermissionID, 10, 64)
		if err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		if err := h.Service.RemovePermission(r.Context(), roleID, permissionID); err != nil {
			panic(err)
		}

		utils.JSON(w, http.StatusOK, map[string]interface{}{
			"message": h.MsgStore.GetMessage(pkgctx.GetLanguage(r.Context()), language.MsgPermissionRemoved),
		})
	})
}

// CreatePermissionModuleAction associates a module action with a permission
func (h *PermissionHandler) CreatePermissionModuleAction(w http.ResponseWriter, r *http.Request) {
	h.HandleValidationError(w, r, func() {
		var req CreatePermissionModuleActionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		if err := h.Service.CreatePermissionModuleAction(r.Context(), req.PermissionID, req.ModuleActionID); err != nil {
			panic(err)
		}

		utils.JSON(w, http.StatusOK, map[string]interface{}{
			"message": h.MsgStore.GetMessage(pkgctx.GetLanguage(r.Context()), language.MsgPermissionModuleActionCreated),
		})
	})
}

// GetPermissionModuleActions returns all module actions for a permission
func (h *PermissionHandler) GetPermissionModuleActions(w http.ResponseWriter, r *http.Request) {
	h.HandleValidationError(w, r, func() {
		permissionID := chi.URLParam(r, "permissionID")
		if permissionID == "" {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		permissionIDInt, err := strconv.ParseInt(permissionID, 10, 64)
		if err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		moduleActions, err := h.Service.GetPermissionModuleActions(r.Context(), permissionIDInt)
		if err != nil {
			panic(err)
		}

		utils.JSON(w, http.StatusOK, moduleActions)
	})
}

// UpdatePermissionModuleActions updates module actions for a permission
func (h *PermissionHandler) UpdatePermissionModuleActions(w http.ResponseWriter, r *http.Request) {
	h.HandleValidationError(w, r, func() {
		var req struct {
			ModuleActionIDs []int64 `json:"module_action_ids" validate:"required"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		permissionID, err := strconv.ParseInt(chi.URLParam(r, "permissionID"), 10, 64)
		if err != nil {
			panic(h.newError(r.Context(), errors.ErrorTypeValidation, errors.ErrorCodeValidationFailed, language.MsgValidationFailed))
		}

		if err := h.Service.UpdatePermissionModuleActions(r.Context(), permissionID, req.ModuleActionIDs); err != nil {
			panic(err)
		}

		utils.JSON(w, http.StatusOK, map[string]interface{}{
			"message": h.MsgStore.GetMessage(pkgctx.GetLanguage(r.Context()), language.MsgPermissionAssigned),
		})
	})
}

// GetModuleActions returns all module actions
func (h *PermissionHandler) GetModuleActions(w http.ResponseWriter, r *http.Request) {
	h.HandleValidationError(w, r, func() {
		actions, err := h.Service.GetModuleActions(r.Context())
		if err != nil {
			panic(err)
		}

		utils.JSON(w, http.StatusOK, actions)
	})
}
