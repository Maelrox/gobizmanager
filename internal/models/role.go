package model

import (
	"time"
)

type Permission struct {
	ID          int64     `json:"id"`
	CompanyID   int64     `json:"company_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Role struct {
	ID          int64        `json:"id"`
	CompanyID   int64        `json:"company_id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type UserRole struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	CompanyUserID int64     `json:"company_user_id"`
	RoleID        int64     `json:"role_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
