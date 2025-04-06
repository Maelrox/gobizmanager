package language

import (
	"sync"
)

// Language codes
const (
	English = "en"
	Spanish = "es"
)

// Message keys
const (
	// Auth messages
	MsgAuthHeaderRequired      = "auth.header_required"
	MsgAuthInvalidFormat       = "auth.invalid_format"
	MsgAuthTokenExpired        = "auth.token_expired"
	MsgAuthInvalidToken        = "auth.invalid_token"
	MsgAuthInvalidCredentials  = "auth.invalid_credentials"
	MsgAuthUserNotFound        = "auth.user_not_found"
	MsgAuthInvalidRequest      = "auth.invalid_request"
	MsgAuthValidationFailed    = "auth.validation_failed"
	MsgAuthUsernameExists      = "auth.username_exists"
	MsgAuthCreateUserFailed    = "auth.create_user_failed"
	MsgAuthTokenGenFailed      = "auth.token_generation_failed"
	MsgAuthUnauthorized        = "auth.unauthorized"
	MsgAuthPermissionDenied    = "auth.permission_denied"
	MsgAuthDatabaseError       = "auth.database_error"
	MsgAuthInvalidRefreshToken = "auth.invalid_refresh_token"
	MsgAuthRegistrationClosed  = "auth.registration_closed"
	MsgAuthInvalidEmail        = "auth.invalid_email"
	MsgAuthPasswordTooShort    = "auth.password_too_short"
	MsgAuthFieldRequired       = "auth.field_required"

	// Rate limit messages
	MsgRateLimitExceeded = "rate_limit.exceeded"

	// Company messages
	MsgCompanyNotFound           = "company.not_found"
	MsgCompanyCreateFailed       = "company.create_failed"
	MsgCompanyUpdateFailed       = "company.update_failed"
	MsgCompanyDeleteFailed       = "company.delete_failed"
	MsgCompanyListFailed         = "company.list_failed"
	MsgCompanyAlreadyExists      = "company.already_exists"
	MsgCompanyNameRequired       = "company.name_required"
	MsgCompanyEmailRequired      = "company.email_required"
	MsgCompanyPhoneRequired      = "company.phone_required"
	MsgCompanyIdentifierRequired = "company.identifier_required"
	MsgCompanyLogoRequired       = "company.logo_required"
	MsgCompanyLogoUpdateFailed   = "company.logo_update_failed"
	MsgCompanyLogoUpdated        = "company.logo_updated"
	MsgCompanyUpdated            = "company.updated"
	MsgCompanyDeleted            = "company.deleted"
	MsgCompanyUserNotFound       = "company.user_not_found"

	// Permission messages
	MsgPermissionDenied       = "permission.denied"
	MsgPermissionCheckFailed  = "permission.check_failed"
	MsgPermissionRequired     = "permission.required"
	MsgPermissionCreateFailed = "permission.create_failed"
	MsgPermissionAssignFailed = "permission.assign_failed"
	MsgPermissionRemoveFailed = "permission.remove_failed"
	MsgPermissionListFailed   = "permission.list_failed"
	MsgRoleCreateFailed       = "role.create_failed"
	MsgRoleNotFound           = "role.not_found"
	MsgRoleListFailed         = "role.list_failed"
	MsgPermissionAssigned     = "permission.assigned"
	MsgPermissionRemoved      = "permission.removed"
	MsgRoleAssignFailed       = "role.assign_failed"
	MsgRoleAssigned           = "role.assigned"

	// Validation messages
	MsgValidationFailed    = "validation.failed"
	MsgValidationRequired  = "validation.required"
	MsgValidationMinLength = "validation.min_length"
	MsgValidationMaxLength = "validation.max_length"
	MsgValidationInvalidID = "validation.invalid_id"

	// Auth validation messages
	MsgValidationRequiredUsername     = "validation.required.username"
	MsgValidationEmailUsername        = "validation.email.username"
	MsgValidationRequiredPassword     = "validation.required.password"
	MsgValidationMinPassword          = "validation.min.password"
	MsgValidationRequiredRefreshToken = "validation.required.refresh_token"

	// Company validation messages
	MsgValidationRequiredName       = "validation.required.name"
	MsgValidationMinName            = "validation.min.name"
	MsgValidationMaxName            = "validation.max.name"
	MsgValidationRequiredEmail      = "validation.required.email"
	MsgValidationEmailEmail         = "validation.email.email"
	MsgValidationRequiredPhone      = "validation.required.phone"
	MsgValidationRequiredIdentifier = "validation.required.identifier"
	MsgValidationRequiredLogo       = "validation.required.logo"
)

// MessageStore holds all language messages
type MessageStore struct {
	mu       sync.RWMutex
	messages map[string]map[string]string
}

// NewMessageStore creates a new message store with default messages
func NewMessageStore() *MessageStore {
	store := &MessageStore{
		messages: make(map[string]map[string]string),
	}

	// Initialize with English messages
	store.messages[English] = map[string]string{
		// Auth messages
		MsgAuthHeaderRequired:      "Authorization header required",
		MsgAuthInvalidFormat:       "Invalid authorization format",
		MsgAuthTokenExpired:        "Token expired",
		MsgAuthInvalidToken:        "Invalid token",
		MsgAuthInvalidCredentials:  "Invalid credentials",
		MsgAuthUserNotFound:        "User not found",
		MsgAuthInvalidRequest:      "Invalid request",
		MsgAuthValidationFailed:    "Validation failed",
		MsgAuthUsernameExists:      "Username already exists",
		MsgAuthCreateUserFailed:    "Failed to create user",
		MsgAuthTokenGenFailed:      "Failed to generate tokens",
		MsgAuthUnauthorized:        "Unauthorized access",
		MsgAuthPermissionDenied:    "Permission denied",
		MsgAuthDatabaseError:       "Database error",
		MsgAuthInvalidRefreshToken: "Invalid refresh token",
		MsgAuthRegistrationClosed:  "Registration is closed. Only the first user can register as ROOT.",
		MsgAuthInvalidEmail:        "Invalid email",
		MsgAuthPasswordTooShort:    "Password too short",
		MsgAuthFieldRequired:       "Field is required",

		// Rate limit messages
		MsgRateLimitExceeded: "Too many requests. Please try again later.",

		// Company messages
		MsgCompanyNotFound:           "Company not found",
		MsgCompanyCreateFailed:       "Failed to create company",
		MsgCompanyUpdateFailed:       "Failed to update company",
		MsgCompanyDeleteFailed:       "Failed to delete company",
		MsgCompanyListFailed:         "Failed to list companies",
		MsgCompanyAlreadyExists:      "Company already exists",
		MsgCompanyNameRequired:       "Company name is required",
		MsgCompanyEmailRequired:      "Company email is required",
		MsgCompanyPhoneRequired:      "Company phone is required",
		MsgCompanyIdentifierRequired: "Company identifier is required",
		MsgCompanyLogoRequired:       "Company logo is required",
		MsgCompanyLogoUpdateFailed:   "Failed to update company logo",
		MsgCompanyLogoUpdated:        "Company logo updated successfully",
		MsgCompanyUpdated:            "Company updated successfully",
		MsgCompanyDeleted:            "Company deleted successfully",
		MsgCompanyUserNotFound:       "Company user not found",

		// Permission messages
		MsgPermissionDenied:       "Insufficient permissions",
		MsgPermissionCheckFailed:  "Failed to check permissions",
		MsgPermissionRequired:     "Permission required",
		MsgPermissionCreateFailed: "Failed to create permission",
		MsgPermissionAssignFailed: "Failed to assign permission to role",
		MsgPermissionRemoveFailed: "Failed to remove permission from role",
		MsgPermissionListFailed:   "Failed to list permissions",
		MsgRoleCreateFailed:       "Failed to create role",
		MsgRoleNotFound:           "Role not found",
		MsgRoleListFailed:         "Failed to list roles",
		MsgPermissionAssigned:     "Permission assigned successfully",
		MsgPermissionRemoved:      "Permission removed successfully",
		MsgRoleAssignFailed:       "Failed to assign role",
		MsgRoleAssigned:           "Role assigned successfully",

		// Validation messages
		MsgValidationFailed:    "Validation failed",
		MsgValidationRequired:  "Field is required",
		MsgValidationMinLength: "Field must be at least %d characters",
		MsgValidationMaxLength: "Field must be at most %d characters",
		MsgValidationInvalidID: "Invalid ID format",

		// Auth validation messages
		MsgValidationRequiredUsername:     "Username is required",
		MsgValidationEmailUsername:        "Username must be a valid email address",
		MsgValidationRequiredPassword:     "Password is required",
		MsgValidationMinPassword:          "Password must be at least 8 characters",
		MsgValidationRequiredRefreshToken: "Refresh token is required",

		// Company validation messages
		MsgValidationRequiredName:       "Company name is required",
		MsgValidationMinName:            "Company name must be at least 3 characters",
		MsgValidationMaxName:            "Company name must be at most 100 characters",
		MsgValidationRequiredEmail:      "Company email is required",
		MsgValidationEmailEmail:         "Company email must be a valid email address",
		MsgValidationRequiredPhone:      "Company phone is required",
		MsgValidationRequiredIdentifier: "Company identifier is required",
		MsgValidationRequiredLogo:       "Company logo is required",
	}

	// Initialize with Spanish messages
	store.messages[Spanish] = map[string]string{
		// Auth messages
		MsgAuthHeaderRequired:      "Se requiere el encabezado de autorización",
		MsgAuthInvalidFormat:       "Formato de autorización inválido",
		MsgAuthTokenExpired:        "Token expirado",
		MsgAuthInvalidToken:        "Token inválido",
		MsgAuthInvalidCredentials:  "Credenciales inválidas",
		MsgAuthUserNotFound:        "Usuario no encontrado",
		MsgAuthInvalidRequest:      "Solicitud inválida",
		MsgAuthValidationFailed:    "Validación fallida",
		MsgAuthUsernameExists:      "El nombre de usuario ya existe",
		MsgAuthCreateUserFailed:    "Error al crear usuario",
		MsgAuthTokenGenFailed:      "Error al generar tokens",
		MsgAuthUnauthorized:        "Acceso no autorizado",
		MsgAuthPermissionDenied:    "Permiso denegado",
		MsgAuthDatabaseError:       "Error de base de datos",
		MsgAuthInvalidRefreshToken: "Token de actualización inválido",
		MsgAuthRegistrationClosed:  "El registro está cerrado. Solo el primer usuario puede registrarse como ROOT.",
		MsgAuthInvalidEmail:        "Correo electrónico inválido",
		MsgAuthPasswordTooShort:    "Contraseña demasiado corta",
		MsgAuthFieldRequired:       "Campo requerido",

		// Rate limit messages
		MsgRateLimitExceeded: "Demasiadas solicitudes. Por favor, intente nuevamente más tarde.",

		// Company messages
		MsgCompanyNotFound:           "Empresa no encontrada",
		MsgCompanyCreateFailed:       "Error al crear la empresa",
		MsgCompanyUpdateFailed:       "Error al actualizar la empresa",
		MsgCompanyDeleteFailed:       "Error al eliminar la empresa",
		MsgCompanyListFailed:         "Error al listar las empresas",
		MsgCompanyAlreadyExists:      "La empresa ya existe",
		MsgCompanyNameRequired:       "Se requiere el nombre de la empresa",
		MsgCompanyEmailRequired:      "Se requiere el correo electrónico de la empresa",
		MsgCompanyPhoneRequired:      "Se requiere el teléfono de la empresa",
		MsgCompanyIdentifierRequired: "Se requiere el identificador de la empresa",
		MsgCompanyLogoRequired:       "Se requiere el logo de la empresa",
		MsgCompanyLogoUpdateFailed:   "Error al actualizar el logo de la empresa",
		MsgCompanyLogoUpdated:        "Logo de la empresa actualizado exitosamente",
		MsgCompanyUpdated:            "Empresa actualizada exitosamente",
		MsgCompanyDeleted:            "Empresa eliminada exitosamente",
		MsgCompanyUserNotFound:       "Usuario de empresa no encontrado",

		// Permission messages
		MsgPermissionDenied:       "Permisos insuficientes",
		MsgPermissionCheckFailed:  "Error al verificar los permisos",
		MsgPermissionRequired:     "Se requieren permisos",
		MsgPermissionCreateFailed: "Error al crear el permiso",
		MsgPermissionAssignFailed: "Error al asignar el permiso al rol",
		MsgPermissionRemoveFailed: "Error al eliminar el permiso del rol",
		MsgPermissionListFailed:   "Error al listar los permisos",
		MsgRoleCreateFailed:       "Error al crear el rol",
		MsgRoleNotFound:           "Rol no encontrado",
		MsgRoleListFailed:         "Error al listar los roles",
		MsgPermissionAssigned:     "Permiso asignado exitosamente",
		MsgPermissionRemoved:      "Permiso eliminado exitosamente",
		MsgRoleAssignFailed:       "Error al asignar el rol",
		MsgRoleAssigned:           "Rol asignado exitosamente",

		// Validation messages
		MsgValidationFailed:    "Error de validación",
		MsgValidationRequired:  "El campo es requerido",
		MsgValidationMinLength: "El campo debe tener al menos %d caracteres",
		MsgValidationMaxLength: "El campo debe tener como máximo %d caracteres",
		MsgValidationInvalidID: "Formato de ID inválido",

		// Auth validation messages
		MsgValidationRequiredUsername:     "El nombre de usuario es requerido",
		MsgValidationEmailUsername:        "El nombre de usuario debe ser un correo electrónico válido",
		MsgValidationRequiredPassword:     "La contraseña es requerida",
		MsgValidationMinPassword:          "La contraseña debe tener al menos 8 caracteres",
		MsgValidationRequiredRefreshToken: "El token de actualización es requerido",

		// Company validation messages
		MsgValidationRequiredName:       "El nombre de la empresa es requerido",
		MsgValidationMinName:            "El nombre de la empresa debe tener al menos 3 caracteres",
		MsgValidationMaxName:            "El nombre de la empresa debe tener como máximo 100 caracteres",
		MsgValidationRequiredEmail:      "El correo electrónico de la empresa es requerido",
		MsgValidationEmailEmail:         "El correo electrónico de la empresa debe ser válido",
		MsgValidationRequiredPhone:      "El teléfono de la empresa es requerido",
		MsgValidationRequiredIdentifier: "El identificador de la empresa es requerido",
		MsgValidationRequiredLogo:       "El logo de la empresa es requerido",
	}

	return store
}

// GetMessage returns the message for the given key in the specified language
func (ms *MessageStore) GetMessage(lang, key string) string {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	// Default to English if language not found
	if _, ok := ms.messages[lang]; !ok {
		lang = English
	}

	if msg, ok := ms.messages[lang][key]; ok {
		return msg
	}

	// Return empty string if message not found
	return ""
}

// AddLanguage adds a new language to the message store
func (ms *MessageStore) AddLanguage(lang string, messages map[string]string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.messages[lang] = messages
}

// UpdateMessage updates a specific message for a language
func (ms *MessageStore) UpdateMessage(lang, key, message string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, ok := ms.messages[lang]; !ok {
		ms.messages[lang] = make(map[string]string)
	}

	ms.messages[lang][key] = message
}
