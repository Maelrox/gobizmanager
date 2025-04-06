package role

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateAdminRole creates an ADMIN role for a company
func (r *Repository) CreateAdminRole(tx *sql.Tx, companyID int64) (int64, error) {
	var adminRoleID int64
	err := tx.QueryRow(`
		INSERT INTO roles (company_id, name, description, created_at, updated_at)
		VALUES (?, 'ADMIN', 'Company administrator', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id
	`, companyID).Scan(&adminRoleID)
	return adminRoleID, err
}

// AssignRoleToUser assigns a role to a user
func (r *Repository) AssignRoleToUser(tx *sql.Tx, userID, roleID int64) error {
	_, err := tx.Exec(
		"INSERT INTO user_roles (user_id, role_id, created_at, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
		userID, roleID,
	)
	return err
}

// GetRoleByID returns a role by ID
func (r *Repository) GetRoleByID(id int64) (*Role, error) {
	var role Role
	err := r.db.QueryRow(`
		SELECT id, company_id, name, description, created_at, updated_at 
		FROM roles WHERE id = ?
	`, id).Scan(
		&role.ID, &role.CompanyID, &role.Name, &role.Description,
		&role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRolesByCompany returns all roles for a company
func (r *Repository) GetRolesByCompany(companyID int64) ([]Role, error) {
	rows, err := r.db.Query(`
		SELECT id, company_id, name, description, created_at, updated_at 
		FROM roles WHERE company_id = ?
	`, companyID)
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
	return roles, nil
}

// DeleteRole deletes a role and its associated permissions
func (r *Repository) DeleteRole(tx *sql.Tx, roleID int64) error {
	// First delete user roles
	_, err := tx.Exec("DELETE FROM user_roles WHERE role_id = ?", roleID)
	if err != nil {
		return err
	}

	// Then delete permissions
	_, err = tx.Exec("DELETE FROM permissions WHERE role_id = ?", roleID)
	if err != nil {
		return err
	}

	// Finally delete the role
	_, err = tx.Exec("DELETE FROM roles WHERE id = ?", roleID)
	return err
}
