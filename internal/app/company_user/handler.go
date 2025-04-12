package company_user

import (
	"errors"
	"net/http"
	"strconv"

	"gobizmanager/internal/app/company"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/logger"
	"gobizmanager/pkg/shared"
	"gobizmanager/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	shared.BaseHandler
	repo        *Repository
	companyRepo *company.Repository
	validator   *validator.Validate
}

func (h *Handler) RegisterCompanyUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.MustGetUserID(w, r)
	if !ok {
		return
	}

	var req RegisterCompanyUserRequest
	if err := utils.ParseRequest(r, &req); err != nil {
		h.RespondError(w, r, err)
		return
	}

	if err := ValidateCreateCompanyUser(&req, userID, h.repo, h.companyRepo, *h.validator, h.MsgStore); err != nil {
		h.RespondError(w, r, err)
		return
	}

	companyUser, err := h.repo.RegisterCompanyUser(&req)
	if err != nil {
		logger.Error(err.Error())
		h.RespondError(w, r, errors.New(language.CompanyCreateFailed))
		return
	}

	utils.JSON(w, http.StatusCreated, companyUser)
}

func (h *Handler) ListCompanyUsers(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "companyID")
	companyIDInt, err := strconv.ParseInt(companyID, 10, 64)
	if err != nil {
		h.RespondError(w, r, errors.New(language.CompanyNotFound))
	}

	users, err := h.repo.ListCompanyUsers(companyIDInt)
	if err != nil {
		logger.Error(err.Error())
		h.RespondError(w, r, errors.New(language.CompanyListFailed))
		return
	}

	utils.JSON(w, http.StatusOK, users)
}

// RemoveCompanyUser removes a user from a company
func (h *Handler) RemoveCompanyUser(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "companyID")
	companyIDInt, err := strconv.ParseInt(companyID, 10, 64)
	if err != nil {
		logger.Error(err.Error())
		h.RespondError(w, r, errors.New(language.CompanyNotFound))
		return
	}

	userID := chi.URLParam(r, "userID")
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		logger.Error(err.Error())
		h.RespondError(w, r, errors.New(language.CompanyUserNotFound))
		return
	}

	//TODO: check permissions

	if err := h.repo.RemoveCompanyUser(companyIDInt, userIDInt); err != nil {
		logger.Error(err.Error())
		h.RespondError(w, r, errors.New(language.CompanyUserRemoveFailed))
		return
	}

	utils.JSON(w, http.StatusNoContent, nil)
}
