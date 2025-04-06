package rbac

import (
	"net/http"

	"gobizmanager/pkg/language"

	"github.com/go-chi/chi/v5"
)

// Routes returns the routes for the RBAC module
func Routes(handler *Handler, msgStore *language.MessageStore) http.Handler {
	r := chi.NewRouter()

	// Permission routes
	r.Route("/permissions", func(r chi.Router) {
		r.Post("/", handler.CreatePermission)
		r.Get("/company/{companyID}", handler.ListPermissions)
		r.Post("/module-actions", handler.CreatePermissionModuleAction)
		r.Delete("/assign", handler.RemovePermission)
	})

	// Role routes
	r.Route("/roles", func(r chi.Router) {
		r.Post("/", handler.CreateRole)
		r.Get("/{id}", handler.GetRole)
		r.Get("/", handler.ListRoles)
		r.Post("/assign", handler.AssignRole)
	})

	return r
}
