package rbac

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	pkgctx "gobizmanager/pkg/context"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/logger"
	"gobizmanager/pkg/utils"

	"go.uber.org/zap"
)

type PermissionHandler struct {
	*RbacBaseHandler
}

func NewPermissionHandler(repo *Repository, msgStore *language.MessageStore) *PermissionHandler {
	return &PermissionHandler{
		RbacBaseHandler: NewBaseHandler(repo, msgStore),
	}
}

func (h *PermissionHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var req CreatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthInvalidRequest))
		return
	}

	permission, err := h.Service.CreatePermission(r.Context(), req.CompanyID, req.Name, req.Description, req.RoleID)
	if err != nil {
		logger.Error("Error creating permission", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.PermissionCreateFailed))
		return
	}

	utils.JSON(w, http.StatusCreated, permission)
}

func (h *PermissionHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "companyID")
	companyIDInt, err := strconv.ParseInt(companyID, 10, 64)
	if err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	permissions, err := h.Service.ListPermissions(r.Context(), companyIDInt)
	if err != nil {
		logger.Error("Error listing permissions", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.PermissionListFailed))
		return
	}

	utils.JSON(w, http.StatusOK, permissions)
}

func (h *PermissionHandler) RemovePermission(w http.ResponseWriter, r *http.Request) {
	var req RemovePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	roleID, err := strconv.ParseInt(req.RoleID, 10, 64)
	if err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	permissionID, err := strconv.ParseInt(req.PermissionID, 10, 64)
	if err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	if err := h.Service.RemovePermission(r.Context(), roleID, permissionID); err != nil {
		logger.Error("Error removing permission", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.PermissionRemoveFailed))
		return
	}

	msg, httpStatus := h.MsgStore.GetMessage(pkgctx.GetLanguage(r.Context()), language.PermissionRemoved)
	utils.JSON(w, httpStatus, msg)
}

// CreatePermissionModuleAction associates a module action with a permission
func (h *PermissionHandler) CreatePermissionModuleAction(w http.ResponseWriter, r *http.Request) {
	var req CreatePermissionModuleActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	if err := h.Service.CreatePermissionModuleAction(r.Context(), req.PermissionID, req.ModuleActionID); err != nil {
		logger.Error("Error creating permission module action", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.PermissionAssignFailed))
		return
	}

	msg, httpStatus := h.MsgStore.GetMessage(pkgctx.GetLanguage(r.Context()), language.PermissionAssigned)
	utils.JSON(w, httpStatus, msg)
}

func (h *PermissionHandler) GetPermissionModuleActions(w http.ResponseWriter, r *http.Request) {
	permissionID := chi.URLParam(r, "permissionID")
	if permissionID == "" {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	permissionIDInt, err := strconv.ParseInt(permissionID, 10, 64)
	if err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	moduleActions, err := h.Service.GetPermissionModuleActions(r.Context(), permissionIDInt)
	if err != nil {
		logger.Error("Error getting permission module actions", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.PermissionListFailed))
		return
	}

	utils.JSON(w, http.StatusOK, moduleActions)
}

func (h *PermissionHandler) UpdatePermissionModuleActions(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ModuleActionIDs []int64 `json:"module_action_ids" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	permissionID, err := strconv.ParseInt(chi.URLParam(r, "permissionID"), 10, 64)
	if err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
		return
	}

	if err := h.Service.UpdatePermissionModuleActions(r.Context(), permissionID, req.ModuleActionIDs); err != nil {
		logger.Error("Error updating permission module actions", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.PermissionAssignFailed))
		return
	}

	msg, httpStatus := h.MsgStore.GetMessage(pkgctx.GetLanguage(r.Context()), language.PermissionAssigned)
	utils.JSON(w, httpStatus, msg)
}

func (h *PermissionHandler) GetModuleActions(w http.ResponseWriter, r *http.Request) {
	actions, err := h.Service.GetModuleActions(r.Context())
	if err != nil {
		logger.Error("Error getting module actions", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.PermissionListFailed))
		return
	}

	utils.JSON(w, http.StatusOK, actions)
}
