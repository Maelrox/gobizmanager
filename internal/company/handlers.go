package company

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"gobizmanager/internal/auth"
	"gobizmanager/internal/permission"
	"gobizmanager/internal/rbac"
	"gobizmanager/internal/role"
	"gobizmanager/internal/user"
	"gobizmanager/pkg/context"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/utils"
)

type Handler struct {
	Repo           *Repository
	RBACRepo       *rbac.Repository
	UserRepo       *user.Repository
	RoleRepo       *role.Repository
	PermissionRepo *permission.Repository
	Validator      *validator.Validate
	MsgStore       *language.MessageStore
}

func NewHandler(
	repo *Repository,
	rbacRepo *rbac.Repository,
	userRepo *user.Repository,
	roleRepo *role.Repository,
	permissionRepo *permission.Repository,
	msgStore *language.MessageStore,
) *Handler {
	return &Handler{
		Repo:           repo,
		RBACRepo:       rbacRepo,
		UserRepo:       userRepo,
		RoleRepo:       roleRepo,
		PermissionRepo: permissionRepo,
		Validator:      validator.New(),
		MsgStore:       msgStore,
	}
}

func (h *Handler) createCompany(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	// Get user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUnauthorized))
		return
	}

	var req CreateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	// Validate the request
	if err := h.Validator.Struct(req); err != nil {
		utils.ValidationError(w, err, lang, h.MsgStore)
		return
	}

	// Check if company with same name already exists for this user
	exists, err := h.Repo.CompanyExistsForUser(userID, req.Name)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgCompanyCreateFailed))
		return
	}
	if exists {
		utils.JSONError(w, http.StatusConflict, h.MsgStore.GetMessage(lang, language.MsgCompanyAlreadyExists))
		return
	}

	// Start transaction
	tx, err := h.Repo.db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgCompanyCreateFailed))
		return
	}
	defer tx.Rollback()

	// Create company
	companyID, err := h.Repo.CreateCompanyWithTx(tx, req.Name, req.Email, req.Phone, req.Address, req.Logo, req.Identifier, userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgCompanyCreateFailed))
		return
	}

	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgCompanyCreateFailed))
		return
	}

	utils.JSON(w, http.StatusCreated, map[string]interface{}{
		"id":         companyID,
		"name":       req.Name,
		"phone":      req.Phone,
		"email":      req.Email,
		"address":    req.Address,
		"logo":       req.Logo,
		"identifier": req.Identifier,
	})
}

func (h *Handler) listCompanies(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUserNotFound))
		return
	}

	// Check if user is ROOT
	isRoot, err := h.UserRepo.IsRoot(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !isRoot {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	// Get companies for this user
	companies, err := h.Repo.ListCompaniesForUser(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgCompanyListFailed))
		return
	}

	utils.JSON(w, http.StatusOK, companies)
}

func (h *Handler) getCompany(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	// Get user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUserNotFound))
		return
	}

	// Check if user is ROOT
	isRoot, err := h.UserRepo.IsRoot(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !isRoot {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	companyID := chi.URLParam(r, "id")
	company, err := h.Repo.GetCompany(companyID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSONError(w, http.StatusNotFound, h.MsgStore.GetMessage(lang, language.MsgCompanyNotFound))
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgCompanyListFailed))
		return
	}

	utils.JSON(w, http.StatusOK, company)
}

func (h *Handler) updateCompany(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	// Get user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUserNotFound))
		return
	}

	// Check if user is ROOT
	isRoot, err := h.UserRepo.IsRoot(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !isRoot {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	// Get company ID from URL
	companyID := chi.URLParam(r, "id")
	if companyID == "" {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationInvalidID))
		return
	}

	// Verify user has access to this company
	hasAccess, err := h.RBACRepo.HasCompanyAccess(userID, companyID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !hasAccess {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	var req UpdateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	// Update company
	if err := h.Repo.UpdateCompany(companyID, req.Name, req.Phone, req.Email, req.Identifier); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgCompanyUpdateFailed))
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"message": h.MsgStore.GetMessage(lang, language.MsgCompanyUpdated),
	})
}

func (h *Handler) deleteCompany(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	// Get user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUserNotFound))
		return
	}

	// Check if user is ROOT
	isRoot, err := h.UserRepo.IsRoot(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !isRoot {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	companyID := chi.URLParam(r, "id")
	companyIDInt, err := strconv.ParseInt(companyID, 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationInvalidID))
		return
	}

	// Start transaction
	tx, err := h.Repo.db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgCompanyDeleteFailed))
		return
	}
	defer tx.Rollback()

	// Delete all associated data
	// 1. Delete company-user relationships
	err = h.RBACRepo.DeleteCompanyUsersWithTx(tx, companyIDInt)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgCompanyDeleteFailed))
		return
	}

	// 2. Delete company roles and permissions
	err = h.RBACRepo.DeleteCompanyRolesWithTx(tx, companyIDInt)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgCompanyDeleteFailed))
		return
	}

	// 3. Delete the company
	err = h.Repo.DeleteCompanyWithTx(tx, companyIDInt)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgCompanyDeleteFailed))
		return
	}

	err = tx.Commit()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgCompanyDeleteFailed))
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"message": h.MsgStore.GetMessage(lang, language.MsgCompanyDeleted),
	})
}

func (h *Handler) updateCompanyLogo(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	// Get user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUserNotFound))
		return
	}

	// Check if user is ROOT
	isRoot, err := h.UserRepo.IsRoot(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !isRoot {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	// Get company ID from URL
	companyID := chi.URLParam(r, "id")
	if companyID == "" {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationInvalidID))
		return
	}

	// Verify user has access to this company
	hasAccess, err := h.RBACRepo.HasCompanyAccess(userID, companyID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgPermissionCheckFailed))
		return
	}
	if !hasAccess {
		utils.JSONError(w, http.StatusForbidden, h.MsgStore.GetMessage(lang, language.MsgPermissionDenied))
		return
	}

	var req UpdateCompanyLogoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgValidationFailed))
		return
	}

	// Update company logo
	if err := h.Repo.UpdateCompanyLogo(companyID, req.Logo); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgCompanyLogoUpdateFailed))
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"message": h.MsgStore.GetMessage(lang, language.MsgCompanyLogoUpdated),
	})
}
