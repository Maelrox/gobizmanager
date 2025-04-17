package company_user

import "time"

type CompanyUser struct {
	ID        int64     `json:"id"`
	CompanyID int64     `json:"company_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterCompanyUserRequest struct {
	CompanyID int64  `json:"company_id" validate:"required"`
	Username  string `json:"username" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
}
