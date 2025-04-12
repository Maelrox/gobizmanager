package user

import (
	"time"

	utils "gobizmanager/pkg/encryption"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email" encrypted:"true"`
	EmailHash string    `json:"-" gorm:"index"`
	Password  string    `json:"-"`
	Phone     string    `json:"phone" encrypted:"true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,email" msg:"auth.invalid_email"`
	Password string `json:"password" validate:"required,min=8" msg:"auth.password_too_short"`
	Phone    string `json:"phone" validate:"required" msg:"auth.field_required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,email" msg:"auth.invalid_email"`
	Password string `json:"password" validate:"required" msg:"auth.field_required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" msg:"auth.field_required"`
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// EncryptSensitiveFields encrypts sensitive fields using the provided key
func (u *User) EncryptSensitiveFields(key string) error {
	var err error
	if u.Email != "" {
		u.Email, err = utils.Encrypt(u.Email, key)
		if err != nil {
			return err
		}
	}
	if u.Phone != "" {
		u.Phone, err = utils.Encrypt(u.Phone, key)
		if err != nil {
			return err
		}
	}
	return nil
}

// DecryptSensitiveFields decrypts sensitive fields using the provided key
func (u *User) DecryptSensitiveFields(key string) error {
	var err error
	if u.Email != "" {
		u.Email, err = utils.Decrypt(u.Email, key)
		if err != nil {
			return err
		}
	}
	if u.Phone != "" {
		u.Phone, err = utils.Decrypt(u.Phone, key)
		if err != nil {
			return err
		}
	}
	return nil
}
