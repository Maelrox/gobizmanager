package permission

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
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
		SELECT p.id, p.role_id, p.module_action_id, p.created_at, p.updated_at,
		       ma.id, ma.module_id, ma.name, ma.description, ma.created_at, ma.updated_at
		FROM permissions p
		JOIN module_actions ma ON p.module_action_id = ma.id
		WHERE p.role_id = ?
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var p Permission
		var ma ModuleAction
		if err := rows.Scan(
			&p.ID, &p.RoleID, &p.ModuleActionID, &p.CreatedAt, &p.UpdatedAt,
			&ma.ID, &ma.ModuleID, &ma.Name, &ma.Description, &ma.CreatedAt, &ma.UpdatedAt,
		); err != nil {
			return nil, err
		}
		p.ModuleAction = &ma
		permissions = append(permissions, p)
	}
	return permissions, nil
}

// DeletePermissionsByRoleID deletes all permissions for a role
func (r *Repository) DeletePermissionsByRoleID(tx *sql.Tx, roleID int64) error {
	_, err := tx.Exec("DELETE FROM permissions WHERE role_id = ?", roleID)
	return err
}
