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

type Role struct {
	ID          int64     `json:"id"`
	CompanyID   int64     `json:"company_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Permission struct {
	ID             int64     `json:"id"`
	RoleID         int64     `json:"role_id"`
	ModuleActionID int64     `json:"module_action_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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

type CreateRoleRequest struct {
	CompanyID   int64  `json:"company_id" validate:"required"`
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"omitempty,max=255"`
}

type CreatePermissionRequest struct {
	RoleID         int64 `json:"role_id" validate:"required"`
	ModuleActionID int64 `json:"module_action_id" validate:"required"`
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
