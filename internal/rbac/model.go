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

// Request/Response types
type CreateCompanyUserRequest struct {
	CompanyID int64 `json:"company_id" validate:"required"`
	UserID    int64 `json:"user_id" validate:"required"`
	IsMain    bool  `json:"is_main"`
}

// CreatePermissionRequest represents the request to create a new permission
type CreatePermissionRequest struct {
	CompanyID   int64  `json:"company_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	RoleID      int64  `json:"role_id" validate:"required"`
}

// CreateRoleRequest represents the request to create a new role
type CreateRoleRequest struct {
	CompanyID   int64  `json:"company_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
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
	UserID int64 `json:"user_id" validate:"required"`
	RoleID int64 `json:"role_id" validate:"required"`
}

type RootGroup struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreatePermissionGroupRequest represents the request to create a permission group
type CreatePermissionGroupRequest struct {
	CompanyID     int64   `json:"company_id" validate:"required" msg:"company.id_required"`
	Name          string  `json:"name" validate:"required,min=3,max=100" msg:"permission.name_required"`
	Description   string  `json:"description" validate:"required" msg:"permission.description_required"`
	PermissionIDs []int64 `json:"permission_ids" validate:"required,min=1" msg:"permission.ids_required"`
}

// CreatePermissionModuleActionRequest represents a request to associate a module action with a permission
type CreatePermissionModuleActionRequest struct {
	PermissionID   int64 `json:"permission_id" validate:"required"`
	ModuleActionID int64 `json:"module_action_id" validate:"required"`
}

// UpdateRolePermissionsRequest represents the request to update role permissions
type UpdateRolePermissionsRequest struct {
	RoleID        string  `json:"role_id" validate:"required"`
	PermissionIDs []int64 `json:"permission_ids" validate:"required"`
}

var moduleActions []struct {
	ID          int64  `json:"id"`
	ModuleName  string `json:"module_name"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Module names
const (
	ModuleCompany = "company"
	ModuleUser    = "user"
	ModuleRole    = "role"
)

// Action names
const (
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

// RolePermission represents the relationship between roles and permissions
type RolePermission struct {
	RoleID       int64 `gorm:"primaryKey"`
	PermissionID int64 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// PermissionModuleAction represents the relationship between permissions and module actions
type PermissionModuleAction struct {
	PermissionID   int64     `json:"permission_id" gorm:"primaryKey"`
	ModuleActionID int64     `json:"module_action_id" gorm:"primaryKey"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
