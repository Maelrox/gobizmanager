package rbac

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(handler *Handler) http.Handler {
	r := chi.NewRouter()

	// Company user routes
	r.Post("/company-users", handler.CreateCompanyUser)

	// Role routes
	r.Group(func(r chi.Router) {
		r.Use(handler.RequirePermission(ModuleRole, ActionCreate))
		r.Post("/roles", handler.CreateRole)
	})

	// Permission routes
	r.Group(func(r chi.Router) {
		r.Use(handler.RequirePermission(ModuleRole, ActionCreate))
		r.Post("/permissions", handler.CreatePermission)
	})

	// User role assignment routes
	r.Group(func(r chi.Router) {
		r.Use(handler.RequirePermission(ModuleRole, ActionUpdate))
		r.Post("/assign-role", handler.AssignRole)
	})

	return r
}
