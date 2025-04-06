package rbac

import (
	"database/sql"
	"fmt"
	"gobizmanager/internal/permission"
	"strconv"
	"time"
)

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

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CompanyUser operations
func (r *Repository) CreateCompanyUser(companyID, userID int64, isMain bool) (int64, error) {
	now := time.Now()
	result, err := r.db.Exec(
		"INSERT INTO company_users (company_id, user_id, is_main, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		companyID, userID, isMain, now, now,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) CreateCompanyUserWithTx(tx *sql.Tx, companyID, userID int64, isMain bool) (int64, error) {
	now := time.Now()
	result, err := tx.Exec(
		"INSERT INTO company_users (company_id, user_id, is_main, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		companyID, userID, isMain, now, now,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) GetCompanyUserByID(id int64) (*CompanyUser, error) {
	var cu CompanyUser
	err := r.db.QueryRow(
		"SELECT id, company_id, user_id, is_main, created_at, updated_at FROM company_users WHERE id = ?",
		id,
	).Scan(&cu.ID, &cu.CompanyID, &cu.UserID, &cu.IsMain, &cu.CreatedAt, &cu.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &cu, nil
}

func (r *Repository) GetCompanyUserByCompanyAndUser(companyID, userID int64) (*CompanyUser, error) {
	var cu CompanyUser
	err := r.db.QueryRow(
		"SELECT id, company_id, user_id, is_main, created_at, updated_at FROM company_users WHERE company_id = ? AND user_id = ?",
		companyID, userID,
	).Scan(&cu.ID, &cu.CompanyID, &cu.UserID, &cu.IsMain, &cu.CreatedAt, &cu.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &cu, nil
}

// Role operations
func (r *Repository) CreateRole(companyID int64, name, description string) (int64, error) {
	now := time.Now()
	result, err := r.db.Exec(
		"INSERT INTO roles (company_id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		companyID, name, description, now, now,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) GetRoleByID(id int64) (*Role, error) {
	var role Role
	err := r.db.QueryRow(
		"SELECT id, company_id, name, description, created_at, updated_at FROM roles WHERE id = ?",
		id,
	).Scan(&role.ID, &role.CompanyID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Permission operations
func (r *Repository) CreatePermission(name, description string, roleID int64) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create permission
	result, err := tx.Exec(`
		INSERT INTO permissions (name, description, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, name, description)
	if err != nil {
		return 0, fmt.Errorf("failed to create permission: %w", err)
	}

	permissionID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get permission ID: %w", err)
	}

	// Associate permission with role
	_, err = tx.Exec(`
		INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, roleID, permissionID)
	if err != nil {
		return 0, fmt.Errorf("failed to associate permission with role: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return permissionID, nil
}

func (r *Repository) GetPermissionsByRoleID(roleID int64) ([]Permission, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.name, p.description, p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
	return permissions, nil
}

// UserRole operations
func (r *Repository) AssignRole(companyUserID, roleID int64) (int64, error) {
	now := time.Now()
	result, err := r.db.Exec(
		"INSERT INTO user_roles (company_user_id, role_id, created_at, updated_at) VALUES (?, ?, ?, ?)",
		companyUserID, roleID, now, now,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) GetUserRoles(companyUserID int64) ([]UserRole, error) {
	rows, err := r.db.Query(
		"SELECT id, company_user_id, role_id, created_at, updated_at FROM user_roles WHERE company_user_id = ?",
		companyUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userRoles []UserRole
	for rows.Next() {
		var ur UserRole
		if err := rows.Scan(&ur.ID, &ur.CompanyUserID, &ur.RoleID, &ur.CreatedAt, &ur.UpdatedAt); err != nil {
			return nil, err
		}
		userRoles = append(userRoles, ur)
	}
	return userRoles, nil
}

// Module operations
func (r *Repository) GetModuleByID(id int64) (*Module, error) {
	var module Module
	err := r.db.QueryRow(
		"SELECT id, name, description, created_at, updated_at FROM modules WHERE id = ?",
		id,
	).Scan(&module.ID, &module.Name, &module.Description, &module.CreatedAt, &module.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &module, nil
}

// ModuleAction operations
func (r *Repository) GetModuleActionByID(id int64) (*ModuleAction, error) {
	var action ModuleAction
	err := r.db.QueryRow(
		"SELECT id, module_id, name, description, created_at, updated_at FROM module_actions WHERE id = ?",
		id,
	).Scan(&action.ID, &action.ModuleID, &action.Name, &action.Description, &action.CreatedAt, &action.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &action, nil
}

// HasPermission checks if a user has a specific permission
func (r *Repository) HasPermission(userID, moduleActionID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM user_roles ur
		JOIN role_permissions rp ON ur.role_id = rp.role_id
		JOIN permission_module_actions pma ON rp.permission_id = pma.permission_id
		WHERE ur.user_id = ? AND pma.module_action_id = ?
	`, userID, moduleActionID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	return count > 0, nil
}

// GetUserPermissions returns all permissions for a user
func (r *Repository) GetUserPermissions(userID int64) ([]permission.Permission, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT p.id, p.name, p.description, p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = ?
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}
	defer rows.Close()

	var permissions []permission.Permission
	for rows.Next() {
		var p permission.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, p)
	}

	return permissions, nil
}

func (r *Repository) DeleteCompanyUsersWithTx(tx *sql.Tx, companyID int64) error {
	_, err := tx.Exec("DELETE FROM company_users WHERE company_id = ?", companyID)
	return err
}

func (r *Repository) DeleteCompanyRolesWithTx(tx *sql.Tx, companyID int64) error {
	// First delete user roles associated with company roles
	_, err := tx.Exec(`
		DELETE ur FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE r.company_id = ?`, companyID)
	if err != nil {
		return err
	}

	// Then delete permissions associated with company roles
	_, err = tx.Exec(`
		DELETE p FROM permissions p
		JOIN roles r ON p.role_id = r.id
		WHERE r.company_id = ?`, companyID)
	if err != nil {
		return err
	}

	// Finally delete the roles
	_, err = tx.Exec("DELETE FROM roles WHERE company_id = ?", companyID)
	return err
}

// GetCompanyUsersByUserID returns all company users for a given user ID
func (r *Repository) GetCompanyUsersByUserID(userID int64) ([]CompanyUser, error) {
	rows, err := r.db.Query(
		"SELECT id, company_id, user_id, is_main, created_at, updated_at FROM company_users WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companyUsers []CompanyUser
	for rows.Next() {
		var cu CompanyUser
		if err := rows.Scan(&cu.ID, &cu.CompanyID, &cu.UserID, &cu.IsMain, &cu.CreatedAt, &cu.UpdatedAt); err != nil {
			return nil, err
		}
		companyUsers = append(companyUsers, cu)
	}
	return companyUsers, nil
}

func (r *Repository) CreateRootGroup(name string) (int64, error) {
	now := time.Now()
	result, err := r.db.Exec(
		"INSERT INTO root_groups (name, created_at, updated_at) VALUES (?, ?, ?)",
		name, now, now,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) GetRootGroupByID(id int64) (*RootGroup, error) {
	var rg RootGroup
	err := r.db.QueryRow(
		"SELECT id, name, created_at, updated_at FROM root_groups WHERE id = ?",
		id,
	).Scan(&rg.ID, &rg.Name, &rg.CreatedAt, &rg.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &rg, nil
}

// HasCompanyAccess checks if a user has access to a company
func (r *Repository) HasCompanyAccess(userID int64, companyID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM company_users
			WHERE user_id = ? AND company_id = ?
		)
	`, userID, companyID).Scan(&exists)
	return exists, err
}

// IsRoot checks if the user has the ROOT role
func (r *Repository) IsRoot(userID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM user_roles ur
			JOIN roles r ON ur.role_id = r.id
			WHERE ur.user_id = ? AND r.name = 'ROOT' AND r.company_id IS NULL
		)
	`
	var isRoot bool
	err := r.db.QueryRow(query, userID).Scan(&isRoot)
	return isRoot, err
}

// CreateRoleWithPermissions creates a new role with permissions
func (r *Repository) CreateRoleWithPermissions(name, description string, permissions []string) (*Role, error) {
	now := time.Now()
	result, err := r.db.Exec(
		"INSERT INTO roles (name, description, created_at, updated_at) VALUES (?, ?, ?, ?)",
		name, description, now, now,
	)
	if err != nil {
		return nil, err
	}

	roleID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	role := &Role{
		ID:          strconv.FormatInt(roleID, 10),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return role, nil
}

// AssignPermissionToRole assigns a permission to a role
func (r *Repository) AssignPermissionToRole(roleID, permissionID string) error {
	now := time.Now()
	_, err := r.db.Exec(
		"INSERT INTO permissions (role_id, module_action_id, created_at, updated_at) VALUES (?, ?, ?, ?)",
		roleID, permissionID, now, now,
	)
	return err
}

// GetRoleWithPermissions returns a role with its permissions
func (r *Repository) GetRoleWithPermissions(roleID string) (*Role, error) {
	// Convert string ID to int64
	id, err := strconv.ParseInt(roleID, 10, 64)
	if err != nil {
		return nil, err
	}

	// Get role details
	role := &Role{}
	err = r.db.QueryRow(`
		SELECT id, company_id, name, description, created_at, updated_at 
		FROM roles WHERE id = ?
	`, id).Scan(
		&role.ID, &role.CompanyID, &role.Name, &role.Description,
		&role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get permissions for the role
	rows, err := r.db.Query(`
		SELECT p.id, p.name, p.description, p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
	role.Permissions = permissions

	return role, nil
}

// ListPermissions returns all permissions for a company
func (r *Repository) ListPermissions(companyID int64) ([]Permission, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT p.id, p.name, p.description, p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN roles r ON rp.role_id = r.id
		WHERE r.company_id = ? OR r.company_id IS NULL
	`, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, p)
	}
	return permissions, nil
}

// ListRolesWithPermissions returns all roles with their permissions
func (r *Repository) ListRolesWithPermissions() ([]Role, error) {
	// Get all roles
	rows, err := r.db.Query(`
		SELECT id, company_id, name, description, created_at, updated_at 
		FROM roles
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(
			&role.ID, &role.CompanyID, &role.Name, &role.Description,
			&role.CreatedAt, &role.UpdatedAt,
		); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	// Get permissions for each role
	for i := range roles {
		id, _ := strconv.ParseInt(roles[i].ID, 10, 64)
		permissions, err := r.GetPermissionsByRoleID(id)
		if err != nil {
			return nil, err
		}
		roles[i].Permissions = permissions
	}

	return roles, nil
}

// RemovePermissionFromRole removes a permission from a role
func (r *Repository) RemovePermissionFromRole(roleID, permissionID string) error {
	_, err := r.db.Exec(`
		DELETE FROM permissions 
		WHERE role_id = ? AND module_action_id = ?
	`, roleID, permissionID)
	return err
}

// GetModuleActionID returns the ID of a module action by module name and action name
func (r *Repository) GetModuleActionID(moduleName, actionName string) (int64, error) {
	var id int64
	err := r.db.QueryRow(`
		SELECT ma.id
		FROM module_actions ma
		JOIN modules m ON ma.module_id = m.id
		WHERE m.name = ? AND ma.name = ?
	`, moduleName, actionName).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to get module action ID: %w", err)
	}
	return id, nil
}
