package rbac

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"gobizmanager/internal/auth"
	pkgctx "gobizmanager/pkg/context"
	"gobizmanager/pkg/errors"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/utils"
)

// BaseHandler contains common functionality for all RBAC handlers
type BaseHandler struct {
	Service  *Service
	MsgStore *language.MessageStore
}

func NewBaseHandler(repo *Repository, msgStore *language.MessageStore) *BaseHandler {
	val := NewValidator(repo, msgStore)
	service := NewService(repo, val)
	return &BaseHandler{
		Service:  service,
		MsgStore: msgStore,
	}
}

// newError creates a new error with the given type, code, and message key
func (h *BaseHandler) newError(ctx context.Context, errorType errors.ErrorType, code errors.ErrorCode, msgKey string) *errors.Error {
	return errors.NewError(errorType, code, h.MsgStore.GetMessage(pkgctx.GetLanguage(ctx), msgKey))
}

// HandleValidationError handles validation errors
func (h *BaseHandler) HandleValidationError(w http.ResponseWriter, r *http.Request, fn func()) {
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(*errors.Error); ok {
				utils.JSON(w, e.GetHTTPStatus(), e)
				return
			}
			utils.JSON(w, http.StatusInternalServerError, h.newError(r.Context(), errors.ErrorTypeInternal, errors.ErrorCodeInternal, language.MsgInternalError))
		}
	}()
	fn()
}

// RequirePermission is a middleware to check if user has permission
func (h *BaseHandler) RequirePermission(moduleName, actionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.HandleValidationError(w, r, func() {
				userID, ok := auth.GetUserID(r.Context())
				if !ok {
					panic(h.newError(r.Context(), errors.ErrorTypeAuthentication, errors.ErrorCodeUserNotFound, language.MsgAuthUserNotFound))
				}

				hasPermission, err := h.Service.CheckPermission(r.Context(), userID, moduleName, actionName)
				if err != nil {
					panic(err)
				}
				if !hasPermission {
					panic(h.newError(r.Context(), errors.ErrorTypeAuthorization, errors.ErrorCodePermissionDenied, language.MsgPermissionDenied))
				}

				next.ServeHTTP(w, r)
			})
		})
	}
}

// CreatePermission creates a new permission
func (h *BaseHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
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

// GetModuleActions returns all module actions
func (h *BaseHandler) GetModuleActions(w http.ResponseWriter, r *http.Request) {
	h.HandleValidationError(w, r, func() {
		actions, err := h.Service.GetModuleActions(r.Context())
		if err != nil {
			panic(err)
		}

		utils.JSON(w, http.StatusOK, actions)
	})
}

// GetPermissionModuleActions returns all module actions for a permission
func (h *BaseHandler) GetPermissionModuleActions(w http.ResponseWriter, r *http.Request) {
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
func (h *BaseHandler) UpdatePermissionModuleActions(w http.ResponseWriter, r *http.Request) {
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
