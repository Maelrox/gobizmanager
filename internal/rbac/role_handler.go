package rbac

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	pkgctx "gobizmanager/pkg/context"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/logger"
	"gobizmanager/pkg/utils"
)

// RoleHandler handles all role-related HTTP requests
type RoleHandler struct {
	*RbacBaseHandler
}

func NewRoleHandler(repo *Repository, msgStore *language.MessageStore) *RoleHandler {
	return &RoleHandler{
		RbacBaseHandler: NewBaseHandler(repo, msgStore),
	}
}

// CreateRole creates a new role
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	role, err := h.Service.CreateRole(r.Context(), req.CompanyID, req.Name, req.Description)
	if err != nil {
		logger.Error("Error creating role", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.RoleCreateFailed))
		return
	}

	utils.JSON(w, http.StatusCreated, role)
}

// GetRole returns a role by ID
func (h *RoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")
	if roleID == "" {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	roleIDInt, err := strconv.ParseInt(roleID, 10, 64)
	if err != nil {
		logger.Error("Error parsing role ID", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	role, err := h.Service.GetRole(r.Context(), roleIDInt)
	if err != nil {
		logger.Error("Error getting role", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.RoleListFailed))
		return
	}

	utils.JSON(w, http.StatusOK, role)
}

// ListRoles returns all roles for a company
func (h *RoleHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "companyID")
	companyIDInt, err := strconv.ParseInt(companyID, 10, 64)
	if err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	roles, err := h.Service.ListRoles(r.Context(), companyIDInt)
	if err != nil {
		logger.Error("Error listing roles", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.RoleListFailed))
		return
	}

	utils.JSON(w, http.StatusOK, roles)
}

// AssignRole assigns a role to a user
func (h *RoleHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	var req AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	if err := h.Service.AssignRole(r.Context(), req.UserID, req.RoleID); err != nil {
		logger.Error("Error assigning role", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.RoleAssignFailed))
	}

	msg, httpStatus := h.MsgStore.GetMessage(pkgctx.GetLanguage(r.Context()), language.RoleAssigned)
	utils.JSON(w, httpStatus, msg)
}

// UpdateRolePermissions updates permissions for a role
func (h *RoleHandler) UpdateRolePermissions(w http.ResponseWriter, r *http.Request) {
	var req UpdateRolePermissionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	roleID, err := strconv.ParseInt(req.RoleID, 10, 64)
	if err != nil {
		logger.Error("Error parsing role ID", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	if err := h.Service.UpdateRolePermissions(r.Context(), roleID, req.PermissionIDs); err != nil {
		logger.Error("Error updating role permissions", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.RoleAssignFailed))
		return
	}

	msg, httpStatus := h.MsgStore.GetMessage(pkgctx.GetLanguage(r.Context()), language.RoleAssigned)
	utils.JSON(w, httpStatus, msg)
}
