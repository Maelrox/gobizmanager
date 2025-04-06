package company

import (
	"database/sql"
	"fmt"

	"gobizmanager/internal/rbac"
	"gobizmanager/platform/config"
)

type Repository struct {
	db       *sql.DB
	cfg      *config.Config
	RBACRepo *rbac.Repository
}

func NewRepository(db *sql.DB, cfg *config.Config, rbacRepo *rbac.Repository) *Repository {
	return &Repository{
		db:       db,
		cfg:      cfg,
		RBACRepo: rbacRepo,
	}
}

func (r *Repository) CreateCompanyWithTx(tx *sql.Tx, name, email, phone, address, logo, identifier string, userID int64) (int64, error) {
	// Create temporary company to handle encryption
	company := &Company{
		Name:       name,
		Email:      email,
		Phone:      phone,
		Address:    address,
		Logo:       sql.NullString{String: logo, Valid: logo != ""},
		Identifier: identifier,
	}

	// Encrypt sensitive fields
	if err := company.EncryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return 0, fmt.Errorf("failed to encrypt company fields: %w", err)
	}

	// Create company
	result, err := tx.Exec(`
		INSERT INTO companies (name, email, phone, address, logo, identifier, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, company.Name, company.Email, company.Phone, company.Address, company.Logo, company.Identifier)
	if err != nil {
		return 0, fmt.Errorf("failed to create company: %w", err)
	}

	companyID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get company ID: %w", err)
	}

	// Create company-user relationship
	_, err = tx.Exec(`
		INSERT INTO company_users (company_id, user_id, is_main, created_at, updated_at)
		VALUES (?, ?, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, companyID, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create company-user relationship: %w", err)
	}

	// Create ADMIN role for the company
	roleResult, err := tx.Exec(`
		INSERT INTO roles (name, description, company_id, created_at, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, "ADMIN", "Company administrator", companyID)
	if err != nil {
		return 0, fmt.Errorf("failed to create ADMIN role: %w", err)
	}

	roleID, err := roleResult.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get role ID: %w", err)
	}

	// Get all permissions
	rows, err := tx.Query("SELECT id FROM permissions")
	if err != nil {
		return 0, fmt.Errorf("failed to get permissions: %w", err)
	}
	defer rows.Close()

	// Assign all permissions to ADMIN role
	for rows.Next() {
		var permissionID int64
		if err := rows.Scan(&permissionID); err != nil {
			return 0, fmt.Errorf("failed to scan permission ID: %w", err)
		}
		_, err = tx.Exec(`
			INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
			VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, roleID, permissionID)
		if err != nil {
			return 0, fmt.Errorf("failed to assign permission to ADMIN role: %w", err)
		}
	}

	// Assign ADMIN role to user
	_, err = tx.Exec(`
		INSERT INTO user_roles (user_id, role_id, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, userID, roleID)
	if err != nil {
		return 0, fmt.Errorf("failed to assign ADMIN role to user: %w", err)
	}

	return companyID, nil
}

func (r *Repository) CreateCompany(name, phone, email, identifier string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO companies (name, phone, email, identifier) 
		VALUES (?, ?, ?, ?)
	`, name, phone, email, identifier)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) GetCompany(id string) (*Company, error) {
	company := &Company{}
	err := r.db.QueryRow(`
		SELECT id, name, phone, email, identifier, logo, created_at, updated_at 
		FROM companies WHERE id = ?
	`, id).Scan(
		&company.ID, &company.Name, &company.Phone, &company.Email,
		&company.Identifier, &company.Logo, &company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Decrypt sensitive fields
	if err := company.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return nil, err
	}

	return company, nil
}

func (r *Repository) UpdateCompany(id string, name, phone, email, identifier string) error {
	// Create a temporary company to encrypt fields
	company := &Company{
		Name:       name,
		Phone:      phone,
		Email:      email,
		Identifier: identifier,
	}

	// Encrypt sensitive fields
	if err := company.EncryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return err
	}

	_, err := r.db.Exec(`
		UPDATE companies 
		SET name = ?, phone = ?, email = ?, identifier = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`, company.Name, company.Phone, company.Email, company.Identifier, id)
	return err
}

func (r *Repository) UpdateCompanyLogo(id string, logo string) error {
	_, err := r.db.Exec(`
		UPDATE companies 
		SET logo = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`, logo, id)
	return err
}

func (r *Repository) DeleteCompanyWithTx(tx *sql.Tx, companyID int64) error {
	// Delete company-user relationships
	if err := r.RBACRepo.DeleteCompanyUsersWithTx(tx, companyID); err != nil {
		return err
	}

	// Delete company roles and permissions
	if err := r.RBACRepo.DeleteCompanyRolesWithTx(tx, companyID); err != nil {
		return err
	}

	// Delete the company
	_, err := tx.Exec("DELETE FROM companies WHERE id = ?", companyID)
	return err
}

func (r *Repository) ListCompanies() ([]Company, error) {
	rows, err := r.db.Query(`
		SELECT id, name, phone, email, identifier, logo, created_at, updated_at 
		FROM companies
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var c Company
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Email,
			&c.Phone,
			&c.Address,
			&c.Identifier,
			&c.Logo,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan company: %w", err)
		}

		// Decrypt sensitive fields
		if err := c.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
			return nil, fmt.Errorf("failed to decrypt company fields: %w", err)
		}

		companies = append(companies, c)
	}
	return companies, nil
}

// ListCompaniesForUser returns all companies for a given user
func (r *Repository) ListCompaniesForUser(userID int64) ([]Company, error) {
	query := `
		SELECT DISTINCT c.id, c.name, c.phone, c.email, c.identifier, c.logo, c.created_at, c.updated_at 
		FROM companies c
		JOIN company_users cu ON c.id = cu.company_id
		WHERE cu.user_id = ?
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var c Company
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Email,
			&c.Phone,
			&c.Address,
			&c.Identifier,
			&c.Logo,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan company: %w", err)
		}

		// Decrypt sensitive fields
		if err := c.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
			return nil, fmt.Errorf("failed to decrypt company fields: %w", err)
		}

		companies = append(companies, c)
	}
	return companies, nil
}

// CompanyExistsForUser checks if a company with the given name already exists for the user
func (r *Repository) CompanyExistsForUser(userID int64, name string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM companies c
			JOIN company_users cu ON c.id = cu.company_id
			WHERE cu.user_id = ? AND c.name = ?
		)`
	err := r.db.QueryRow(query, userID, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check company existence: %w", err)
	}
	return exists, nil
}

// AddUserToCompany adds a user to a company
func (r *Repository) AddUserToCompany(tx *sql.Tx, companyID, userID int64) error {
	query := `INSERT INTO company_users (company_id, user_id, created_at, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	_, err := tx.Exec(query, companyID, userID)
	if err != nil {
		return fmt.Errorf("failed to add user to company: %w", err)
	}
	return nil
}
