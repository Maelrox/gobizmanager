package user

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// SearchUsers handles listing users within a company
func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	companyID := r.URL.Query().Get("companyId")
	if companyID == "" {
		http.Error(w, "Company ID is required", http.StatusBadRequest)
		return
	}

	users, err := h.repo.SearchUsers(companyID)
	if err != nil {
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
