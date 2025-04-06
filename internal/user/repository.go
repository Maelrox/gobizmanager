// internal/user/repository.go
package user

import (
	"database/sql"

	"gobizmanager/pkg/encryption"
	"gobizmanager/platform/config"
	"time"
)

type Repository struct {
	db  *sql.DB
	cfg *config.Config
}

func NewRepository(db *sql.DB, cfg *config.Config) *Repository {
	return &Repository{db: db, cfg: cfg}
}

// CreateUserWithTx creates a new user within a transaction
func (r *Repository) CreateUserWithTx(tx *sql.Tx, username, password string) (int64, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	result, err := tx.Exec(
		"INSERT INTO users (email, password, created_at, updated_at) VALUES (?, ?, ?, ?)",
		username, hashedPassword, now, now,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetUserByID returns a user by ID
func (r *Repository) GetUserByID(id int64) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(`
		SELECT id, email, password, phone, created_at, updated_at 
		FROM users WHERE id = ?
	`, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Phone,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Decrypt sensitive fields
	if err := user.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail returns a user by email
func (r *Repository) GetUserByEmail(email string) (*User, error) {
	// Encrypt the email for comparison
	encryptedEmail, err := encryption.Encrypt(email, r.cfg.EncryptionKey)
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = r.db.QueryRow(`
		SELECT id, email, password, phone, created_at, updated_at 
		FROM users WHERE email = ?
	`, encryptedEmail).Scan(
		&user.ID, &user.Email, &user.Password, &user.Phone,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Decrypt sensitive fields
	if err := user.DecryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return nil, err
	}

	return user, nil
}

// GetDB returns the database connection
func (r *Repository) GetDB() *sql.DB {
	return r.db
}

func (r *Repository) CreateUser(email, password, phone string) (int64, error) {
	// Create a temporary user to encrypt fields
	user := &User{
		Email:    email,
		Password: password,
		Phone:    phone,
	}

	// Encrypt sensitive fields
	if err := user.EncryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return 0, err
	}

	result, err := r.db.Exec(`
		INSERT INTO users (email, password, phone) 
		VALUES (?, ?, ?)
	`, user.Email, user.Password, user.Phone)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) UpdateUser(id int64, email, password, phone string) error {
	// Create a temporary user to encrypt fields
	user := &User{
		Email:    email,
		Password: password,
		Phone:    phone,
	}

	// Encrypt sensitive fields
	if err := user.EncryptSensitiveFields(r.cfg.EncryptionKey); err != nil {
		return err
	}

	_, err := r.db.Exec(`
		UPDATE users 
		SET email = ?, password = ?, phone = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`, user.Email, user.Password, user.Phone, id)
	return err
}

// GetRootRoleID returns the ID of the ROOT role
func (r *Repository) GetRootRoleID(tx *sql.Tx) (int64, error) {
	var rootRoleID int64
	err := tx.QueryRow("SELECT id FROM roles WHERE name = 'ROOT' AND company_id IS NULL").Scan(&rootRoleID)
	if err != nil {
		return 0, err
	}
	return rootRoleID, nil
}

// AssignRootRole assigns the ROOT role to a user
func (r *Repository) AssignRootRole(tx *sql.Tx, userID, roleID int64) error {
	_, err := tx.Exec(
		"INSERT INTO user_roles (user_id, role_id, created_at, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
		userID, roleID,
	)
	return err
}

// RegisterRootUser registers a new ROOT user with all necessary operations in a transaction
func (r *Repository) RegisterRootUser(username, password string) (int64, error) {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // This ensures rollback if we don't commit

	// Get ROOT role ID
	rootRoleID, err := r.GetRootRoleID(tx)
	if err != nil {
		return 0, err
	}

	// Create user
	userID, err := r.CreateUserWithTx(tx, username, password)
	if err != nil {
		return 0, err
	}

	// Assign ROOT role
	if err := r.AssignRootRole(tx, userID, rootRoleID); err != nil {
		return 0, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return userID, nil
}

// RegisterUser handles the complete user registration process
func (r *Repository) RegisterUser(username, password string) (int64, error) {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Create user
	userID, err := r.CreateUserWithTx(tx, username, password)
	if err != nil {
		return 0, err
	}

	// Get ROOT role ID
	rootRoleID, err := r.GetRootRoleID(tx)
	if err != nil {
		return 0, err
	}

	// Assign ROOT role
	if err := r.AssignRootRole(tx, userID, rootRoleID); err != nil {
		return 0, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return userID, nil
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
