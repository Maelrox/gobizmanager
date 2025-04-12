package company_user

import (
	"errors"

	"gobizmanager/internal/app/company"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/logger"

	"github.com/go-playground/validator/v10"
)

func ValidateCreateCompanyUser(req *RegisterCompanyUserRequest, userID int64, repo *Repository, companyRepo *company.Repository, validator validator.Validate, msgStore *language.MessageStore) error {
	if err := validator.Struct(req); err != nil {
		return err
	}

	exists, err := companyRepo.CompanyExistsForUserByID(userID, req.CompanyID)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(language.CompanyListFailed)
	}
	if !exists {
		return errors.New(language.CompanyNotFound)
	}
	return nil
}
