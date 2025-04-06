package rbac

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(handler *Handler) http.Handler {
	r := chi.NewRouter()

	// Role routes
	r.Route("/roles", func(r chi.Router) {
		r.Post("/", handler.CreateRole)
		r.Get("/", handler.ListRoles)
		r.Get("/{id}", handler.GetRole)
	})

	// Permission routes
	r.Route("/permissions", func(r chi.Router) {
		r.Post("/{companyID}", handler.CreatePermission)
		r.Get("/{companyID}", handler.ListPermissions)
	})

	// Role-Permission association routes
	r.Route("/role-permissions", func(r chi.Router) {
		r.Post("/assign", handler.AssignPermission)
		r.Post("/remove", handler.RemovePermission)
	})

	// User role assignment routes
	r.Group(func(r chi.Router) {
		r.Use(handler.RequirePermission(ModuleRole, ActionUpdate))
		r.Post("/assign-role", handler.AssignRole)
	})

	return r
}
