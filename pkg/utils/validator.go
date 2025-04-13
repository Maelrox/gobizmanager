package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	pkgctx "gobizmanager/pkg/context"
	"gobizmanager/pkg/language"

	"github.com/go-playground/validator/v10"
)

// GetValidationMessage extracts the validation message from the struct tag
func GetValidationMessage(err validator.FieldError, lang string, msgStore *language.MessageStore) string {
	// First try to get the message from the validation tag
	fieldName := err.Field()
	validationTag := err.Tag()

	// Try to get a custom message based on the validation tag and field name
	msgKey := fmt.Sprintf("validation.%s.%s", validationTag, strings.ToLower(fieldName))
	if msg, _ := msgStore.GetMessage(lang, msgKey); msg != "" {
		return msg
	}

	// If no custom message, use default messages
	switch validationTag {
	case "required":
		return fmt.Sprintf("%s is required", strings.ToLower(fieldName))
	case "min", "max":
		paramStr := err.Param()
		param, err := strconv.Atoi(paramStr)
		if err != nil {
			msg, _ := msgStore.GetMessage(lang, language.ValidationFailed)
			return msg
		}
		if validationTag == "min" {
			return fmt.Sprintf("%s must be at least %d characters", strings.ToLower(fieldName), param)
		}
		return fmt.Sprintf("%s must be at most %d characters", strings.ToLower(fieldName), param)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", strings.ToLower(fieldName))
	default:
		msg, _ := msgStore.GetMessage(lang, language.ValidationFailed)
		return msg
	}
}

func ValidationError(w http.ResponseWriter, r *http.Request, err error, msgStore *language.MessageStore) {
	lang := pkgctx.GetLanguage(r.Context())
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errorMsg := GetValidationMessage(validationErrors[0], lang, msgStore)
		JSONError(w, http.StatusBadRequest, errorMsg)
		return
	}
	msg, httpStatus := msgStore.GetMessage(lang, language.ValidationFailed)
	JSONError(w, httpStatus, msg)
}

func GetValidationError(err error, lang string, msgStore *language.MessageStore) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errorMsg := GetValidationMessage(validationErrors[0], lang, msgStore)
		return errorMsg
	}
	return "undefined error"
}

func ParseRequest(r *http.Request, req interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.New(language.BadRequest)
	}
	return nil
}
