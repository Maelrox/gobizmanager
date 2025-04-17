package validation

import (
	"fmt"
	"strings"
	"unicode"
)

type PasswordValidator struct {
	MinLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireNumbers   bool
	RequireSpecial   bool
}

func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		MinLength:        8,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSpecial:   true,
	}
}

func (v *PasswordValidator) Validate(password, username string) error {
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

	if strings.Contains(strings.ToLower(password), strings.ToLower(username)) {
		return fmt.Errorf("password cannot contain username")
	}

	if hasKeyboardPattern(password) {
		return fmt.Errorf("password contains a common keyboard pattern")
	}

	if hasRepeatedChars(password, 3) {
		return fmt.Errorf("password contains too many repeated characters")
	}

	if hasSequentialNumbers(password) {
		return fmt.Errorf("password contains sequential numbers")
	}

	if hasSequentialLetters(password) {
		return fmt.Errorf("password contains sequential letters")
	}

	return nil
}

// hasKeyboardPattern checks for common keyboard patterns like qwerty, asdf, etc.
func hasKeyboardPattern(password string) bool {
	commonPatterns := []string{
		"qwerty", "asdfgh", "zxcvbn", "qwertz", "azerty", "123456", "123456789", "1234567890", "12345678901234567890",
	}

	lower := strings.ToLower(password)
	for _, pattern := range commonPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}

func hasRepeatedChars(password string, threshold int) bool {
	if len(password) < threshold {
		return false
	}

	count := 1
	prev := rune(password[0])

	for _, char := range password[1:] {
		if char == prev {
			count++
			if count >= threshold {
				return true
			}
		} else {
			count = 1
			prev = char
		}
	}
	return false
}

// hasSequentialNumbers checks for sequences like "123", "456", etc.
func hasSequentialNumbers(password string) bool {
	for i := 0; i < len(password)-2; i++ {
		if unicode.IsDigit(rune(password[i])) &&
			unicode.IsDigit(rune(password[i+1])) &&
			unicode.IsDigit(rune(password[i+2])) {
			n1 := int(password[i] - '0')
			n2 := int(password[i+1] - '0')
			n3 := int(password[i+2] - '0')

			// Check ascending or descending sequence
			if (n2 == n1+1 && n3 == n2+1) || (n2 == n1-1 && n3 == n2-1) {
				return true
			}
		}
	}
	return false
}

// hasSequentialLetters checks for sequences like "abc", "xyz", etc.
func hasSequentialLetters(password string) bool {
	lower := strings.ToLower(password)
	for i := 0; i < len(lower)-2; i++ {
		if unicode.IsLetter(rune(lower[i])) &&
			unicode.IsLetter(rune(lower[i+1])) &&
			unicode.IsLetter(rune(lower[i+2])) {
			if lower[i+1] == lower[i]+1 && lower[i+2] == lower[i+1]+1 {
				return true
			}
		}
	}
	return false
}
