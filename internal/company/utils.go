package company

import (
	"database/sql"
	"gobizmanager/internal/rbac"
	"gobizmanager/internal/role"
	"gobizmanager/internal/role/permission"
)

// CreateCompanyWithAdmin creates a new company and sets up the admin role and permissions
func CreateCompanyWithAdmin(
	tx *sql.Tx,
	repo *Repository,
	rbacRepo *rbac.Repository,
	roleRepo *role.Repository,
	permissionRepo *permission.Repository,
	name, phone, email, identifier string,
	userID int64,
) (int64, error) {
	// Create company
	companyID, err := repo.CreateCompanyWithTx(tx, name, email, phone, "", "", identifier, userID)
	if err != nil {
		return 0, err
	}

	// Create company-user relationship
	_, err = rbacRepo.CreateCompanyUserWithTx(tx, companyID, userID, true)
	if err != nil {
		return 0, err
	}

	// Create ADMIN role for this company
	adminRoleID, err := roleRepo.CreateAdminRole(tx, companyID)
	if err != nil {
		return 0, err
	}

	// Create default permissions for the company
	defaultPermissions := []struct {
		name        string
		description string
	}{
		{"manage_companies", "Full access to company management"},
		{"manage_users", "Full access to user management"},
		{"manage_roles", "Full access to role management"},
	}

	for _, perm := range defaultPermissions {
		_, err := rbacRepo.CreatePermission(companyID, perm.name, perm.description, adminRoleID)
		if err != nil {
			return 0, err
		}
	}

	// Assign ADMIN role to user
	if err := roleRepo.AssignRoleToUser(tx, userID, adminRoleID); err != nil {
		return 0, err
	}

	return companyID, nil
}
