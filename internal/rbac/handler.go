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

			companyIDStr := chi.URLParam(r, "companyID")
			if companyIDStr == "" {
				utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, "company.company_id_required"))
				return
			}

			companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
			if err != nil {
				utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, "company.invalid_company_id"))
				return
			}

			hasPermission, err := h.Repo.HasPermission(userID, companyID, moduleName, actionName)
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

// Create company user
func (h *Handler) CreateCompanyUser(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	var req CreateCompanyUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, "rbac.invalid_request"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, "rbac.validation_failed"))
		return
	}

	// Check if user is already associated with the company
	_, err := h.Repo.GetCompanyUserByCompanyAndUser(req.CompanyID, req.UserID)
	if err == nil {
		utils.JSONError(w, http.StatusConflict, h.MsgStore.GetMessage(lang, "rbac.user_already_associated"))
		return
	}

	id, err := h.Repo.CreateCompanyUser(req.CompanyID, req.UserID, req.IsMain)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, "rbac.create_company_user_failed"))
		return
	}

	utils.JSON(w, http.StatusCreated, map[string]int64{"id": id})
}

// Create role
func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	var req CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, "rbac.invalid_request"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, "rbac.validation_failed"))
		return
	}

	id, err := h.Repo.CreateRole(req.CompanyID, req.Name, req.Description)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, "rbac.create_role_failed"))
		return
	}

	utils.JSON(w, http.StatusCreated, map[string]int64{"id": id})
}

// Create permission
func (h *Handler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	var req CreatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, "rbac.invalid_request"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, "rbac.validation_failed"))
		return
	}

	id, err := h.Repo.CreatePermission(req.RoleID, req.ModuleActionID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, "rbac.create_permission_failed"))
		return
	}

	utils.JSON(w, http.StatusCreated, map[string]int64{"id": id})
}

// Assign role to user
func (h *Handler) AssignRole(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	var req AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, "rbac.invalid_request"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, "rbac.validation_failed"))
		return
	}

	id, err := h.Repo.AssignRole(req.CompanyUserID, req.RoleID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, "rbac.assign_role_failed"))
		return
	}

	utils.JSON(w, http.StatusCreated, map[string]int64{"id": id})
}
