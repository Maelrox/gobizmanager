package rbac

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"

	"gobizmanager/internal/auth"
	pkgctx "gobizmanager/pkg/context"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/logger"
	"gobizmanager/pkg/utils"
)

type RbacBaseHandler struct {
	Service  *Service
	MsgStore *language.MessageStore
}

func NewBaseHandler(repo *Repository, msgStore *language.MessageStore) *RbacBaseHandler {
	val := NewValidator(repo, msgStore)
	service := NewService(repo, val)
	return &RbacBaseHandler{
		Service:  service,
		MsgStore: msgStore,
	}
}

// RequirePermission is a middleware to check if user has permission to access a resource
func (h *RbacBaseHandler) RequirePermission(moduleName, actionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := auth.GetUserID(r.Context())
			if !ok {
				utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthInvalidRequest))
				return
			}

			hasPermission, err := h.Service.CheckPermission(r.Context(), userID, moduleName, actionName)
			if err != nil {
				logger.Error("Error checking permission", zap.Error(err))
				utils.RespondError(w, r, h.MsgStore, errors.New(language.PermissionCheckFailed))
				return
			}
			if !hasPermission {
				utils.RespondError(w, r, h.MsgStore, errors.New(language.PermissionDenied))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (h *RbacBaseHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var req CreatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
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

func (h *RbacBaseHandler) GetModuleActions(w http.ResponseWriter, r *http.Request) {
	actions, err := h.Service.GetModuleActions(r.Context())
	if err != nil {
		logger.Error("Error getting module actions", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.PermissionListFailed))
		return
	}

	utils.JSON(w, http.StatusOK, actions)
}

func (h *RbacBaseHandler) GetPermissionModuleActions(w http.ResponseWriter, r *http.Request) {
	permissionID := chi.URLParam(r, "permissionID")
	if permissionID == "" {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.ValidationFailed))
	}

	permissionIDInt, err := strconv.ParseInt(permissionID, 10, 64)
	if err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.PermissionNotFound))
		return
	}

	moduleActions, err := h.Service.GetPermissionModuleActions(r.Context(), permissionIDInt)
	if err != nil {
		panic(err)
	}

	utils.JSON(w, http.StatusOK, moduleActions)
}

func (h *RbacBaseHandler) UpdatePermissionModuleActions(w http.ResponseWriter, r *http.Request) {
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
