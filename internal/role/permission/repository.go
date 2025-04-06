package permission

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreatePermission creates a new permission
func (r *Repository) CreatePermission(name, description string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO permissions (name, description, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, name, description)
	if err != nil {
		return 0, fmt.Errorf("failed to create permission: %w", err)
	}

	return result.LastInsertId()
}

// AddModuleActionToPermission adds a module action to a permission
func (r *Repository) AddModuleActionToPermission(permissionID, moduleActionID int64) error {
	_, err := r.db.Exec(`
		INSERT INTO permission_module_actions (permission_id, module_action_id, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, permissionID, moduleActionID)
	if err != nil {
		return fmt.Errorf("failed to add module action to permission: %w", err)
	}

	return nil
}

// AssignPermissionToRole assigns a permission to a role
func (r *Repository) AssignPermissionToRole(roleID, permissionID int64) error {
	_, err := r.db.Exec(`
		INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to assign permission to role: %w", err)
	}

	return nil
}

// GetRolePermissions returns all permissions for a role
func (r *Repository) GetRolePermissions(roleID int64) ([]Permission, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.name, p.description, p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
	`, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
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

// GetPermissionModuleActions returns all module actions for a permission
func (r *Repository) GetPermissionModuleActions(permissionID int64) ([]int64, error) {
	rows, err := r.db.Query(`
		SELECT module_action_id
		FROM permission_module_actions
		WHERE permission_id = ?
	`, permissionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission module actions: %w", err)
	}
	defer rows.Close()

	var moduleActionIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan module action ID: %w", err)
		}
		moduleActionIDs = append(moduleActionIDs, id)
	}

	return moduleActionIDs, nil
}

// GrantAllPermissions grants all module actions to a role
func (r *Repository) GrantAllPermissions(tx *sql.Tx, roleID int64) error {
	_, err := tx.Exec(`
		INSERT INTO permissions (role_id, module_action_id, created_at, updated_at)
		SELECT ?, id, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		FROM module_actions
	`, roleID)
	return err
}

// GetPermissionsByRoleID returns all permissions for a role
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

// DeletePermissionsByRoleID deletes all permissions for a role
func (r *Repository) DeletePermissionsByRoleID(tx *sql.Tx, roleID int64) error {
	_, err := tx.Exec("DELETE FROM permissions WHERE role_id = ?", roleID)
	return err
}
