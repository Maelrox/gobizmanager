package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	if msg := msgStore.GetMessage(lang, msgKey); msg != "" {
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
			return msgStore.GetMessage(lang, language.MsgValidationFailed)
		}
		if validationTag == "min" {
			return fmt.Sprintf("%s must be at least %d characters", strings.ToLower(fieldName), param)
		}
		return fmt.Sprintf("%s must be at most %d characters", strings.ToLower(fieldName), param)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", strings.ToLower(fieldName))
	default:
		return msgStore.GetMessage(lang, language.MsgValidationFailed)
	}
}

// ValidationError sends a validation error response
func ValidationError(w http.ResponseWriter, err error, lang string, msgStore *language.MessageStore) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// Get the first error message
		errorMsg := GetValidationMessage(validationErrors[0], lang, msgStore)
		JSONError(w, http.StatusBadRequest, errorMsg)
		return
	}
	JSONError(w, http.StatusBadRequest, msgStore.GetMessage(lang, language.MsgValidationFailed))
}
