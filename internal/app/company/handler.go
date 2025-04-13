package company

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"gobizmanager/internal/app/rbac"
	"gobizmanager/internal/app/user"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/logger"
	"gobizmanager/pkg/shared"
	"gobizmanager/pkg/utils"
)

type Handler struct {
	shared.BaseHandler
	repo      *Repository
	rbacRepo  *rbac.Repository
	userRepo  *user.Repository
	Validator *validator.Validate
}

func (h *Handler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.MustGetUserID(w, r)
	if !ok {
		return
	}
	var req CreateCompanyRequest
	if err := utils.ParseRequest(r, &req); err != nil {
		h.RespondError(w, r, err)
		return
	}

	if err := ValidateCreateCompany(r, &req, userID, h.repo, *h.Validator, h.MsgStore); err != nil {
		h.RespondError(w, r, err)
		return
	}

	company, err := h.repo.CreateCompany(&req, userID)
	if err != nil {
		logger.Error(err.Error())
		h.RespondError(w, r, errors.New(language.CompanyCreateFailed))
		return
	}

	res := CompanyResponse{
		CompanyID:  company.ID,
		Name:       company.Name,
		Phone:      company.Phone,
		Email:      company.Email,
		Address:    company.Address,
		Logo:       "",
		Identifier: company.Identifier,
	}
	utils.JSON(w, http.StatusCreated, res)
}

func (h *Handler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.MustGetUserID(w, r)
	if !ok {
		return
	}

	companies, err := h.repo.ListCompanies(userID)
	if err != nil {
		logger.Error(err.Error())
		h.RespondError(w, r, errors.New(language.CompanyListFailed))
		return
	}

	utils.JSON(w, http.StatusOK, companies)
}

func (h *Handler) GetCompany(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "companyID")
	companyIDInt, err := strconv.ParseInt(companyID, 10, 64)
	if err != nil {
		h.RespondError(w, r, err)
		return
	}

	company, err := h.repo.GetCompany(companyIDInt)
	if err != nil {
		logger.Error(err.Error())
		h.RespondError(w, r, errors.New(language.CompanyGetFailed))
		return
	}

	utils.JSON(w, http.StatusOK, company)
}

func (h *Handler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.MustGetUserID(w, r)
	if !ok {
		return
	}
	companyID := chi.URLParam(r, "companyID")
	companyIDInt, err := strconv.ParseInt(companyID, 10, 64)
	if err != nil {
		h.RespondError(w, r, errors.New(language.CompanyGetFailed))
		return
	}

	var req UpdateCompanyRequest
	if err := utils.ParseRequest(r, &req); err != nil {
		h.RespondError(w, r, err)
		return
	}

	if err := ValidateUpdateCompany(r, &req, userID, h.repo, *h.Validator, h.MsgStore); err != nil {
		h.RespondError(w, r, err)
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		h.RespondError(w, r, err)
		return
	}

	company, err := h.repo.UpdateCompany(companyIDInt, &req)
	if err != nil {
		logger.Error(err.Error())
		h.RespondError(w, r, errors.New(language.CompanyUpdateFailed))
		return
	}
	res := CompanyResponse{
		CompanyID:  company.ID,
		Name:       company.Name,
		Phone:      company.Phone,
		Email:      company.Email,
		Address:    company.Address,
		Logo:       "",
		Identifier: company.Identifier,
	}
	utils.JSON(w, http.StatusOK, res)
}

func (h *Handler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "companyID")
	companyIDInt, err := strconv.ParseInt(companyID, 10, 64)
	if err != nil {
		h.RespondError(w, r, err)
		return
	}

	if err := h.repo.DeleteCompany(companyIDInt); err != nil {
		logger.Error(err.Error())
		h.RespondError(w, r, errors.New(language.CompanyDeleteFailed))
		return
	}

	utils.JSON(w, http.StatusNoContent, nil)
}

func (h *Handler) UpdateCompanyLogo(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}
