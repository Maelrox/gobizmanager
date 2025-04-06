package company_user

import "time"

// CompanyUser represents a user associated with a company
type CompanyUser struct {
	ID        int64     `json:"id"`
	CompanyID int64     `json:"company_id"`
	UserID    int64     `json:"user_id"`
	IsMain    bool      `json:"is_main"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterCompanyUserRequest represents the request to register a new user for a company
type RegisterCompanyUserRequest struct {
	CompanyID int64  `json:"company_id" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	IsMain    bool   `json:"is_main"`
}
