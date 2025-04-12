package rbac

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns the routes for the RBAC module
func Routes(roleHandler *RoleHandler, permissionHandler *PermissionHandler) http.Handler {
	r := chi.NewRouter()

	// Module actions route
	r.Get("/module-actions", permissionHandler.GetModuleActions)

	// Permission routes
	r.Route("/permissions", func(r chi.Router) {
		r.Post("/", permissionHandler.CreatePermission)
		r.Get("/company/{companyID}", permissionHandler.ListPermissions)
		r.Post("/module-actions", permissionHandler.CreatePermissionModuleAction)
		r.Delete("/assign", permissionHandler.RemovePermission)
		r.Get("/{permissionID}/module-actions", permissionHandler.GetPermissionModuleActions)
		r.Put("/{permissionID}/module-actions", permissionHandler.UpdatePermissionModuleActions)
	})

	// Role routes
	r.Route("/roles", func(r chi.Router) {
		r.Post("/", roleHandler.CreateRole)
		r.Get("/{id}", roleHandler.GetRole)
		r.Get("/company/{companyID}", roleHandler.ListRoles)
		r.Post("/assign", roleHandler.AssignRole)
		r.Put("/permissions", roleHandler.UpdateRolePermissions)
	})

	return r
}
