package company_user

import "net/http"

func Routes(handler *Handler) http.Handler {
	r := http.NewServeMux()

	// Register company user routes
	r.HandleFunc("/register/company-user", handler.RegisterCompanyUser)

	return r
}
