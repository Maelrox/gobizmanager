package validation

import (
	"fmt"
	"unicode"
)

type PasswordValidator struct {
	MinLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireNumbers   bool
	RequireSpecial   bool
	DisallowCommon   bool
	commonPasswords  map[string]struct{} // Pre-loaded common passwords
}

func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		MinLength:        8,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSpecial:   true,
		DisallowCommon:   true,
		commonPasswords:  initCommonPasswords(),
	}
}

func (v *PasswordValidator) Validate(password string) error {
	if len(password) < v.MinLength {
		return fmt.Errorf("password must be at least %d characters long", v.MinLength)
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if v.RequireUppercase && !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if v.RequireLowercase && !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if v.RequireNumbers && !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if v.RequireSpecial && !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	if v.DisallowCommon {
		if _, exists := v.commonPasswords[password]; exists {
			return fmt.Errorf("password is too common")
		}
	}

	return nil
}

func initCommonPasswords() map[string]struct{} {
	common := map[string]struct{}{
		"password":   {},
		"123456":     {},
		"12345678":   {},
		"1234567890": {},
	}
	return common
}
