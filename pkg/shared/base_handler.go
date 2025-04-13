package shared

import (
	"errors"
	"net/http"

	"gobizmanager/internal/auth"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/utils"
)

type BaseHandler struct {
	MsgStore *language.MessageStore
}

func (h *BaseHandler) RespondError(w http.ResponseWriter, r *http.Request, err error) {
	utils.RespondError(w, r, h.MsgStore, err)
}
func (h *BaseHandler) MustGetUserID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		h.RespondError(w, r, errors.New(language.AuthUnauthorized))
	}
	return userID, ok
}
