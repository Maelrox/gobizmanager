package company

import (
	"database/sql"
	"gobizmanager/internal/permission"
	"gobizmanager/internal/rbac"
	"gobizmanager/internal/role"
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
	companyID, err := repo.CreateCompanyWithTx(tx, name, phone, email, identifier)
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

	// Grant all permissions to the ADMIN role
	if err := permissionRepo.GrantAllPermissions(tx, adminRoleID); err != nil {
		return 0, err
	}

	// Assign ADMIN role to the user for this company
	if err := roleRepo.AssignRoleToUser(tx, userID, adminRoleID); err != nil {
		return 0, err
	}

	return companyID, nil
}
