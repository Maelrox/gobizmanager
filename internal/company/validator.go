package company

import (
	"errors"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/logger"
	"gobizmanager/pkg/utils"
	"net/http"

	pkgctx "gobizmanager/pkg/context"

	"github.com/go-playground/validator/v10"
)

func ValidateCreateCompany(r *http.Request, req *CreateCompanyRequest, userID int64, repo *Repository, validator validator.Validate, msgStore *language.MessageStore) error {
	lang := pkgctx.GetLanguage(r.Context())

	if err := validator.Struct(req); err != nil {
		return errors.New(utils.GetValidationError(err, lang, msgStore))
	}

	exists, err := repo.CompanyExistsForUser(userID, req.Name)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(language.CompanyListFailed)
	}
	if exists {
		return errors.New(language.CompanyAlreadyExists)
	}
	return nil
}

func ValidateUpdateCompany(r *http.Request, req *UpdateCompanyRequest, userID int64, repo *Repository, validator validator.Validate, msgStore *language.MessageStore) error {
	lang := pkgctx.GetLanguage(r.Context())

	if err := validator.Struct(req); err != nil {
		return errors.New(utils.GetValidationError(err, lang, msgStore))
	}

	if err := validator.Struct(req); err != nil {
		return err
	}

	exists, err := repo.CompanyExistsForUser(userID, req.Name)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(language.CompanyListFailed)
	}
	if !exists {
		return errors.New(language.CompanyNotFound)
	}
	return nil
}
