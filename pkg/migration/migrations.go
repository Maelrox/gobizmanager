package migration

import (
	"database/sql"
	"fmt"
)

var migrations = []struct {
	name string
	stmt string
}{
	{
		name: "Create companies table",
		stmt: `CREATE TABLE IF NOT EXISTS companies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT,
			phone TEXT,
			address TEXT,
			identifier TEXT,
			logo TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	},
	{
		name: "Create users table",
		stmt: `CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			email_hash TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			phone TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	},
	{
		name: "Create company_users table",
		stmt: `CREATE TABLE IF NOT EXISTS company_users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			company_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			is_main BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(company_id, user_id)
		)`,
	},
	{
		name: "Create modules table",
		stmt: `CREATE TABLE IF NOT EXISTS modules (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			description TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,
	},
	{
		name: "Create module_actions table",
		stmt: `CREATE TABLE IF NOT EXISTS module_actions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			module_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE CASCADE,
			UNIQUE(module_id, name)
		)`,
	},
	{
		name: "Create roles table",
		stmt: `CREATE TABLE IF NOT EXISTS roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			company_id INTEGER,
			name TEXT NOT NULL,
			description TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
			UNIQUE(company_id, name)
		)`,
	},
	{
		name: "Create permissions table",
		stmt: `CREATE TABLE IF NOT EXISTS permissions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(name)
		)`,
	},
	{
		name: "Create permission_module_actions table",
		stmt: `CREATE TABLE IF NOT EXISTS permission_module_actions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			permission_id INTEGER NOT NULL,
			module_action_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
			FOREIGN KEY (module_action_id) REFERENCES module_actions(id) ON DELETE CASCADE,
			UNIQUE(permission_id, module_action_id)
		)`,
	},
	{
		name: "Create role_permissions table",
		stmt: `CREATE TABLE IF NOT EXISTS role_permissions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			role_id INTEGER NOT NULL,
			permission_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
			FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
			UNIQUE(role_id, permission_id)
		)`,
	},
	{
		name: "Create user_roles table",
		stmt: `CREATE TABLE IF NOT EXISTS user_roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			role_id INTEGER NOT NULL,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
			UNIQUE(user_id, role_id)
		)`,
	},
	{
		name: "Create default modules and roles",
		stmt: `
			-- Create default modules
			INSERT OR IGNORE INTO modules (id, name, description, created_at, updated_at)
			VALUES 
				(1, 'company', 'Company management module', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(2, 'user', 'User management module', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(3, 'role', 'Role management module', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

			-- Create module actions
			INSERT OR IGNORE INTO module_actions (id, module_id, name, description, created_at, updated_at)
			VALUES 
				(1, 1, 'create', 'Create company', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(2, 1, 'read', 'View company', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(3, 1, 'update', 'Update company', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(4, 1, 'delete', 'Delete company', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(5, 2, 'create', 'Create user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(6, 2, 'read', 'View user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(7, 2, 'update', 'Update user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(8, 2, 'delete', 'Delete user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(9, 3, 'create', 'Create role', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(10, 3, 'read', 'View role', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(11, 3, 'update', 'Update role', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(12, 3, 'delete', 'Delete role', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

			-- Create ROOT role (with NULL company_id for system-wide access)
			INSERT OR IGNORE INTO roles (id, company_id, name, description, created_at, updated_at)
			VALUES (1, NULL, 'ROOT', 'System ROOT user with full access', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

			-- Create default permissions
			INSERT OR IGNORE INTO permissions (id, name, description, created_at, updated_at)
			VALUES 
				(1, 'manage_companies', 'Full access to company management', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(2, 'manage_users', 'Full access to user management', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(3, 'manage_roles', 'Full access to role management', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

			-- Add module actions to permissions
			INSERT OR IGNORE INTO permission_module_actions (permission_id, module_action_id, created_at, updated_at)
			SELECT 1, id, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
			FROM module_actions
			WHERE module_id = 1;

			INSERT OR IGNORE INTO permission_module_actions (permission_id, module_action_id, created_at, updated_at)
			SELECT 2, id, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
			FROM module_actions
			WHERE module_id = 2;

			INSERT OR IGNORE INTO permission_module_actions (permission_id, module_action_id, created_at, updated_at)
			SELECT 3, id, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
			FROM module_actions
			WHERE module_id = 3;

			-- Grant all permissions to ROOT role
			INSERT OR IGNORE INTO role_permissions (role_id, permission_id, created_at, updated_at)
			VALUES 
				(1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(1, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
				(1, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
		`,
	},
}

func ApplyMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Apply each migration if not already applied
	for _, m := range migrations {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", m.name).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if count == 0 {
			tx, err := db.Begin()
			if err != nil {
				return fmt.Errorf("failed to begin transaction: %w", err)
			}

			_, err = tx.Exec(m.stmt)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to apply migration '%s': %w", m.name, err)
			}

			_, err = tx.Exec("INSERT INTO migrations (name) VALUES (?)", m.name)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to record migration '%s': %w", m.name, err)
			}

			if err := tx.Commit(); err != nil {
				return fmt.Errorf("failed to commit transaction: %w", err)
			}

			fmt.Printf("Applied migration: %s\n", m.name)
		}
	}

	return nil
}
