package rbac

import "time"

type CompanyUser struct {
	ID        int64     `json:"id"`
	CompanyID int64     `json:"company_id"`
	UserID    int64     `json:"user_id"`
	IsMain    bool      `json:"is_main"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Module struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ModuleAction struct {
	ID          int64     `json:"id"`
	ModuleID    int64     `json:"module_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Permission represents a system permission
type Permission struct {
	ID             string    `json:"id"`
	RoleID         int64     `json:"role_id"`
	ModuleActionID int64     `json:"module_action_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Role represents a system role
type Role struct {
	ID          string       `json:"id"`
	CompanyID   int64        `json:"company_id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type UserRole struct {
	ID            int64     `json:"id"`
	CompanyUserID int64     `json:"company_user_id"`
	RoleID        int64     `json:"role_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Request/Response types
type CreateCompanyUserRequest struct {
	CompanyID int64 `json:"company_id" validate:"required"`
	UserID    int64 `json:"user_id" validate:"required"`
	IsMain    bool  `json:"is_main"`
}

// CreatePermissionRequest represents the request to create a new permission
type CreatePermissionRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// CreateRoleRequest represents the request to create a new role
type CreateRoleRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Permissions []string `json:"permissions" validate:"required"`
}

// AssignPermissionRequest represents the request to assign a permission to a role
type AssignPermissionRequest struct {
	RoleID       string `json:"role_id" validate:"required"`
	PermissionID string `json:"permission_id" validate:"required"`
}

// RemovePermissionRequest represents the request to remove a permission from a role
type RemovePermissionRequest struct {
	RoleID       string `json:"role_id" validate:"required"`
	PermissionID string `json:"permission_id" validate:"required"`
}

type AssignRoleRequest struct {
	CompanyUserID int64 `json:"company_user_id" validate:"required"`
	RoleID        int64 `json:"role_id" validate:"required"`
}

type RootGroup struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
