package company_user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"gobizmanager/internal/auth"
	"gobizmanager/internal/rbac"
	"gobizmanager/pkg/context"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/utils"
)

type Handler struct {
	Repo      *Repository
	RBACRepo  *rbac.Repository
	Validator *validator.Validate
	MsgStore  *language.MessageStore
}

func NewHandler(repo *Repository, rbacRepo *rbac.Repository, msgStore *language.MessageStore) *Handler {
	return &Handler{
		Repo:      repo,
		RBACRepo:  rbacRepo,
		Validator: validator.New(),
		MsgStore:  msgStore,
	}
}

// RegisterCompanyUser handles the registration of a new user for a company
func (h *Handler) RegisterCompanyUser(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	// Get user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, "auth.unauthorized"))
		return
	}

	// Check if user has manage_users permission
	moduleActionID, err := h.RBACRepo.GetModuleActionID("user", "manage")
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, "rbac.permission_check_failed"))
		return
	}

	hasPermission, err := h.RBACRepo.HasPermission(userID, moduleActionID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, "rbac.permission_check_failed"))
		return
	}
	if !hasPermission {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, "rbac.insufficient_permissions"))
		return
	}

	var req RegisterCompanyUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, "validation.invalid_request"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, "validation.failed"))
		return
	}

	// Verify user has access to the company
	companyIDStr := strconv.FormatInt(req.CompanyID, 10)
	hasAccess, err := h.RBACRepo.HasCompanyAccess(userID, companyIDStr)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, "rbac.permission_check_failed"))
		return
	}
	if !hasAccess {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, "rbac.insufficient_permissions"))
		return
	}

	companyUser, err := h.Repo.RegisterCompanyUser(&req)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, "company_user.registration_failed"))
		return
	}

	utils.JSON(w, http.StatusCreated, companyUser)
}
