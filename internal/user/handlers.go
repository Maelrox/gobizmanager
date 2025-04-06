package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// SearchUsers handles searching for users within a company for dropdown
func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "companyId")
	query := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "50"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}
	if limit > 50 {
		limit = 50
	}

	users, err := h.repo.SearchUsers(companyID, query, limit)
	if err != nil {
		http.Error(w, "Failed to search users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
