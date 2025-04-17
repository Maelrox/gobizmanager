package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	pkgctx "gobizmanager/pkg/context"
	"gobizmanager/pkg/language"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func JSONError(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}

func RespondError(w http.ResponseWriter, r *http.Request, msgStore *language.MessageStore, err error) {
	msg, httpStatus := msgStore.GetMessage(pkgctx.GetLanguage(r.Context()), err.Error())
	JSONError(w, httpStatus, msg)
}

func ParseRequest(r *http.Request, req interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.New(language.BadRequest)
	}
	return nil
}
