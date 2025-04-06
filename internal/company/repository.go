package company

import (
	"database/sql"

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

func (r *Repository) CreateCompanyWithTx(tx *sql.Tx, name, phone, email, identifier string) (int64, error) {
	// Create a temporary company to encrypt fields
	company := &Company{
		Name:       name,
		Phone:      phone,
		Email:      email,
		Identifier: identifier,
	}

	// Encrypt sensitive fields
	if err := company.EncryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return 0, err
	}

	result, err := tx.Exec(`
		INSERT INTO companies (name, phone, email, identifier) 
		VALUES (?, ?, ?, ?)
	`, company.Name, company.Phone, company.Email, company.Identifier)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
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
			&c.ID, &c.Name, &c.Phone, &c.Email,
			&c.Identifier, &c.Logo, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// Decrypt sensitive fields
		if err := c.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
			return nil, err
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
			&c.ID, &c.Name, &c.Phone, &c.Email,
			&c.Identifier, &c.Logo, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// Decrypt sensitive fields
		if err := c.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
			return nil, err
		}

		companies = append(companies, c)
	}
	return companies, nil
}
