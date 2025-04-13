package language

import (
	"errors"
	"net/http"
	"sync"
)

const (
	English = "en"
	Spanish = "es"
)

// Message keys
const (
	// Shared validation messages
	BadRequest = "bad.request"

	// Auth messages
	AuthHeaderRequired      = "auth.header_required"
	AuthInvalidFormat       = "auth.invalid_format"
	AuthTokenExpired        = "auth.token_expired"
	AuthInvalidToken        = "auth.invalid_token"
	AuthInvalidCredentials  = "auth.invalid_credentials"
	AuthUserNotFound        = "auth.user_not_found"
	AuthInvalidRequest      = "auth.invalid_request"
	AuthValidationFailed    = "auth.validation_failed"
	AuthUsernameExists      = "auth.username_exists"
	AuthCreateUserFailed    = "auth.create_user_failed"
	AuthTokenGenFailed      = "auth.token_generation_failed"
	AuthUnauthorized        = "auth.unauthorized"
	AuthPermissionDenied    = "auth.permission_denied"
	AuthDatabaseError       = "auth.database_error"
	AuthInvalidRefreshToken = "auth.invalid_refresh_token"
	AuthRegistrationClosed  = "auth.registration_closed"
	AuthInvalidEmail        = "auth.invalid_email"
	AuthPasswordTooShort    = "auth.password_too_short"
	AuthFieldRequired       = "auth.field_required"

	// Rate limit messages
	RateLimitExceeded = "rate_limit.exceeded"

	// Company messages
	CompanyNotFound           = "company.not_found"
	CompanyCreateFailed       = "company.create_failed"
	CompanyUpdateFailed       = "company.update_failed"
	CompanyDeleteFailed       = "company.delete_failed"
	CompanyListFailed         = "company.list_failed"
	CompanyGetFailed          = "company.get_failed"
	CompanyAlreadyExists      = "company.already_exists"
	CompanyNameRequired       = "company.name_required"
	CompanyEmailRequired      = "company.email_required"
	CompanyPhoneRequired      = "company.phone_required"
	CompanyIdentifierRequired = "company.identifier_required"
	CompanyLogoRequired       = "company.logo_required"
	CompanyLogoUpdateFailed   = "company.logo_update_failed"
	CompanyLogoUpdated        = "company.logo_updated"
	CompanyUpdated            = "company.updated"
	CompanyDeleted            = "company.deleted"
	CompanyUserNotFound       = "company.user_not_found"
	CompanyUserRemoveFailed   = "company.user_remove_failed"

	// Permission messages
	PermissionDenied       = "permission.denied"
	PermissionCheckFailed  = "permission.check_failed"
	PermissionRequired     = "permission.required"
	PermissionCreateFailed = "permission.create_failed"
	PermissionAssignFailed = "permission.assign_failed"
	PermissionRemoveFailed = "permission.remove_failed"
	PermissionListFailed   = "permission.list_failed"
	RoleCreateFailed       = "role.create_failed"
	RoleNotFound           = "role.not_found"
	RoleListFailed         = "role.list_failed"
	PermissionAssigned     = "permission.assigned"
	PermissionRemoved      = "permission.removed"
	RoleAssignFailed       = "role.assign_failed"
	RoleAssigned           = "role.assigned"
	PermissionNotFound     = "permission.not_found"

	// Module actions messages
	ModuleActionCreated = "module.action.denied"

	// Validation messages
	ValidationFailed    = "validation.failed"
	ValidationRequired  = "validation.required"
	ValidationMinLength = "validation.min_length"
	ValidationMaxLength = "validation.max_length"
	ValidationInvalidID = "validation.invalid_id"

	// Auth validation messages
	ValidationRequiredUsername     = "validation.required.username"
	ValidationEmailUsername        = "validation.email.username"
	ValidationRequiredPassword     = "validation.required.password"
	ValidationMinPassword          = "validation.min.password"
	ValidationRequiredRefreshToken = "validation.required.refresh_token"

	// Company validation messages
	ValidationRequiredName       = "validation.required.name"
	ValidationMinName            = "validation.min.name"
	ValidationMaxName            = "validation.max.name"
	ValidationRequiredEmail      = "validation.required.email"
	ValidationEmailEmail         = "validation.email.email"
	ValidationRequiredPhone      = "validation.required.phone"
	ValidationRequiredIdentifier = "validation.required.identifier"
	ValidationRequiredLogo       = "validation.required.logo"
)

// Message represents a localized message with its HTTP status code
type Message struct {
	Text   string
	Status int
}

// MessageStore holds all language messages
type MessageStore struct {
	mu       sync.RWMutex
	messages map[string]map[string]Message
}

func (m *MessageStore) New(code string) error {
	return errors.New(code) // later can be enhanced to return custom error types
}

// NewMessageStore creates a new message store with default messages
func NewMessageStore() *MessageStore {
	store := &MessageStore{
		messages: make(map[string]map[string]Message),
	}

	// Initialize with English messages
	store.messages[English] = map[string]Message{
		// Shared validation messages
		BadRequest: {"Bad request", http.StatusBadRequest},

		// Auth messages
		AuthHeaderRequired:      {"Authorization header required", http.StatusUnauthorized},
		AuthInvalidFormat:       {"Invalid authorization format", http.StatusBadRequest},
		AuthTokenExpired:        {"Token expired", http.StatusUnauthorized},
		AuthInvalidToken:        {"Invalid token", http.StatusUnauthorized},
		AuthInvalidCredentials:  {"Invalid credentials", http.StatusUnauthorized},
		AuthUserNotFound:        {"User not found", http.StatusNotFound},
		AuthInvalidRequest:      {"Invalid request", http.StatusBadRequest},
		AuthValidationFailed:    {"Validation failed", http.StatusBadRequest},
		AuthUsernameExists:      {"Username already exists", http.StatusConflict},
		AuthCreateUserFailed:    {"Failed to create user", http.StatusInternalServerError},
		AuthTokenGenFailed:      {"Failed to generate tokens", http.StatusInternalServerError},
		AuthUnauthorized:        {"Unauthorized access", http.StatusUnauthorized},
		AuthPermissionDenied:    {"Permission denied", http.StatusForbidden},
		AuthDatabaseError:       {"Database error", http.StatusInternalServerError},
		AuthInvalidRefreshToken: {"Invalid refresh token", http.StatusUnauthorized},
		AuthRegistrationClosed:  {"Registration is closed. Only the first user can register as ROOT.", http.StatusForbidden},
		AuthInvalidEmail:        {"Invalid email", http.StatusBadRequest},
		AuthPasswordTooShort:    {"Password too short", http.StatusBadRequest},
		AuthFieldRequired:       {"Field is required", http.StatusBadRequest},

		// Rate limit messages
		RateLimitExceeded: {"Too many requests. Please try again later.", http.StatusTooManyRequests},

		// Company messages
		CompanyNotFound:           {"Company not found", http.StatusNotFound},
		CompanyCreateFailed:       {"Failed to create company", http.StatusInternalServerError},
		CompanyUpdateFailed:       {"Failed to update company", http.StatusInternalServerError},
		CompanyDeleteFailed:       {"Failed to delete company", http.StatusInternalServerError},
		CompanyListFailed:         {"Failed to list companies", http.StatusInternalServerError},
		CompanyGetFailed:          {"Failed to get company", http.StatusInternalServerError},
		CompanyAlreadyExists:      {"Company already exists", http.StatusConflict},
		CompanyNameRequired:       {"Company name is required", http.StatusBadRequest},
		CompanyEmailRequired:      {"Company email is required", http.StatusBadRequest},
		CompanyPhoneRequired:      {"Company phone is required", http.StatusBadRequest},
		CompanyIdentifierRequired: {"Company identifier is required", http.StatusBadRequest},
		CompanyLogoRequired:       {"Company logo is required", http.StatusBadRequest},
		CompanyLogoUpdateFailed:   {"Failed to update company logo", http.StatusInternalServerError},
		CompanyLogoUpdated:        {"Company logo updated successfully", http.StatusOK},
		CompanyUpdated:            {"Company updated successfully", http.StatusOK},
		CompanyDeleted:            {"Company deleted successfully", http.StatusOK},
		CompanyUserNotFound:       {"Company user not found", http.StatusNotFound},
		CompanyUserRemoveFailed:   {"Failed to remove company user", http.StatusInternalServerError},

		// Permission messages
		PermissionDenied:       {"Insufficient permissions", http.StatusForbidden},
		PermissionCheckFailed:  {"Failed to check permissions", http.StatusInternalServerError},
		PermissionRequired:     {"Permission required", http.StatusForbidden},
		PermissionCreateFailed: {"Failed to create permission", http.StatusInternalServerError},
		PermissionAssignFailed: {"Failed to assign permission to role", http.StatusInternalServerError},
		PermissionRemoveFailed: {"Failed to remove permission from role", http.StatusInternalServerError},
		PermissionListFailed:   {"Failed to list permissions", http.StatusInternalServerError},
		RoleCreateFailed:       {"Failed to create role", http.StatusInternalServerError},
		RoleNotFound:           {"Role not found", http.StatusNotFound},
		RoleListFailed:         {"Failed to list roles", http.StatusInternalServerError},
		PermissionAssigned:     {"Permission assigned successfully", http.StatusOK},
		PermissionRemoved:      {"Permission removed successfully", http.StatusOK},
		RoleAssignFailed:       {"Failed to assign role", http.StatusInternalServerError},
		RoleAssigned:           {"Role assigned successfully", http.StatusOK},
		PermissionNotFound:     {"Permission not found", http.StatusNotFound},

		// Validation messages
		ValidationFailed:    {"Validation failed", http.StatusBadRequest},
		ValidationRequired:  {"Field is required", http.StatusBadRequest},
		ValidationMinLength: {"Field must be at least %d characters", http.StatusBadRequest},
		ValidationMaxLength: {"Field must be at most %d characters", http.StatusBadRequest},
		ValidationInvalidID: {"Invalid ID format", http.StatusBadRequest},

		// Auth validation messages
		ValidationRequiredUsername:     {"Username is required", http.StatusBadRequest},
		ValidationEmailUsername:        {"Username must be a valid email address", http.StatusBadRequest},
		ValidationRequiredPassword:     {"Password is required", http.StatusBadRequest},
		ValidationMinPassword:          {"Password must be at least 8 characters", http.StatusBadRequest},
		ValidationRequiredRefreshToken: {"Refresh token is required", http.StatusBadRequest},

		// Company validation messages
		ValidationRequiredName:       {"Company name is required", http.StatusBadRequest},
		ValidationMinName:            {"Company name must be at least 3 characters", http.StatusBadRequest},
		ValidationMaxName:            {"Company name must be at most 100 characters", http.StatusBadRequest},
		ValidationRequiredEmail:      {"Company email is required", http.StatusBadRequest},
		ValidationEmailEmail:         {"Company email must be a valid email address", http.StatusBadRequest},
		ValidationRequiredPhone:      {"Company phone is required", http.StatusBadRequest},
		ValidationRequiredIdentifier: {"Company identifier is required", http.StatusBadRequest},
		ValidationRequiredLogo:       {"Company logo is required", http.StatusBadRequest},
	}

	// Initialize with Spanish messages
	store.messages[Spanish] = map[string]Message{
		// Shared validation messages
		BadRequest: {"Solicitud inválida", http.StatusBadRequest},

		// Auth messages
		AuthHeaderRequired:      {"Se requiere el encabezado de autorización", http.StatusUnauthorized},
		AuthInvalidFormat:       {"Formato de autorización inválido", http.StatusBadRequest},
		AuthTokenExpired:        {"Token expirado", http.StatusUnauthorized},
		AuthInvalidToken:        {"Token inválido", http.StatusUnauthorized},
		AuthInvalidCredentials:  {"Credenciales inválidas", http.StatusUnauthorized},
		AuthUserNotFound:        {"Usuario no encontrado", http.StatusNotFound},
		AuthInvalidRequest:      {"Solicitud inválida", http.StatusBadRequest},
		AuthValidationFailed:    {"Validación fallida", http.StatusBadRequest},
		AuthUsernameExists:      {"El nombre de usuario ya existe", http.StatusConflict},
		AuthCreateUserFailed:    {"Error al crear usuario", http.StatusInternalServerError},
		AuthTokenGenFailed:      {"Error al generar tokens", http.StatusInternalServerError},
		AuthUnauthorized:        {"Acceso no autorizado", http.StatusUnauthorized},
		AuthPermissionDenied:    {"Permiso denegado", http.StatusForbidden},
		AuthDatabaseError:       {"Error de base de datos", http.StatusInternalServerError},
		AuthInvalidRefreshToken: {"Token de actualización inválido", http.StatusUnauthorized},
		AuthRegistrationClosed:  {"El registro está cerrado. Solo el primer usuario puede registrarse como ROOT.", http.StatusForbidden},
		AuthInvalidEmail:        {"Correo electrónico inválido", http.StatusBadRequest},
		AuthPasswordTooShort:    {"Contraseña demasiado corta", http.StatusBadRequest},
		AuthFieldRequired:       {"Campo requerido", http.StatusBadRequest},

		// Rate limit messages
		RateLimitExceeded: {"Demasiadas solicitudes. Por favor, intente nuevamente más tarde.", http.StatusTooManyRequests},

		// Company messages
		CompanyNotFound:           {"Empresa no encontrada", http.StatusNotFound},
		CompanyCreateFailed:       {"Error al crear la empresa", http.StatusInternalServerError},
		CompanyUpdateFailed:       {"Error al actualizar la empresa", http.StatusInternalServerError},
		CompanyDeleteFailed:       {"Error al eliminar la empresa", http.StatusInternalServerError},
		CompanyListFailed:         {"Error al listar las empresas", http.StatusInternalServerError},
		CompanyGetFailed:          {"Error al obtener la empresa", http.StatusInternalServerError},
		CompanyAlreadyExists:      {"La empresa ya existe", http.StatusConflict},
		CompanyNameRequired:       {"Se requiere el nombre de la empresa", http.StatusBadRequest},
		CompanyEmailRequired:      {"Se requiere el correo electrónico de la empresa", http.StatusBadRequest},
		CompanyPhoneRequired:      {"Se requiere el teléfono de la empresa", http.StatusBadRequest},
		CompanyIdentifierRequired: {"Se requiere el identificador de la empresa", http.StatusBadRequest},
		CompanyLogoRequired:       {"Se requiere el logo de la empresa", http.StatusBadRequest},
		CompanyLogoUpdateFailed:   {"Error al actualizar el logo de la empresa", http.StatusInternalServerError},
		CompanyLogoUpdated:        {"Logo de la empresa actualizado exitosamente", http.StatusOK},
		CompanyUpdated:            {"Empresa actualizada exitosamente", http.StatusOK},
		CompanyDeleted:            {"Empresa eliminada exitosamente", http.StatusOK},
		CompanyUserNotFound:       {"Usuario de empresa no encontrado", http.StatusNotFound},
		CompanyUserRemoveFailed:   {"Error al eliminar el usuario de la empresa", http.StatusInternalServerError},

		// Permission messages
		PermissionDenied:       {"Permisos insuficientes", http.StatusForbidden},
		PermissionCheckFailed:  {"Error al verificar los permisos", http.StatusInternalServerError},
		PermissionRequired:     {"Se requieren permisos", http.StatusForbidden},
		PermissionCreateFailed: {"Error al crear el permiso", http.StatusInternalServerError},
		PermissionAssignFailed: {"Error al asignar el permiso al rol", http.StatusInternalServerError},
		PermissionRemoveFailed: {"Error al eliminar el permiso del rol", http.StatusInternalServerError},
		PermissionListFailed:   {"Error al listar los permisos", http.StatusInternalServerError},
		RoleCreateFailed:       {"Error al crear el rol", http.StatusInternalServerError},
		RoleNotFound:           {"Rol no encontrado", http.StatusNotFound},
		RoleListFailed:         {"Error al listar los roles", http.StatusInternalServerError},
		PermissionAssigned:     {"Permiso asignado exitosamente", http.StatusOK},
		PermissionRemoved:      {"Permiso eliminado exitosamente", http.StatusOK},
		RoleAssignFailed:       {"Error al asignar el rol", http.StatusInternalServerError},
		RoleAssigned:           {"Rol asignado exitosamente", http.StatusOK},
		PermissionNotFound:     {"Permiso no encontrado", http.StatusNotFound},

		// Validation messages
		ValidationFailed:    {"Error de validación", http.StatusBadRequest},
		ValidationRequired:  {"El campo es requerido", http.StatusBadRequest},
		ValidationMinLength: {"El campo debe tener al menos %d caracteres", http.StatusBadRequest},
		ValidationMaxLength: {"El campo debe tener como máximo %d caracteres", http.StatusBadRequest},
		ValidationInvalidID: {"Formato de ID inválido", http.StatusBadRequest},

		// Auth validation messages
		ValidationRequiredUsername:     {"El nombre de usuario es requerido", http.StatusBadRequest},
		ValidationEmailUsername:        {"El nombre de usuario debe ser un correo electrónico válido", http.StatusBadRequest},
		ValidationRequiredPassword:     {"La contraseña es requerida", http.StatusBadRequest},
		ValidationMinPassword:          {"La contraseña debe tener al menos 8 caracteres", http.StatusBadRequest},
		ValidationRequiredRefreshToken: {"El token de actualización es requerido", http.StatusBadRequest},

		// Company validation messages
		ValidationRequiredName:       {"El nombre de la empresa es requerido", http.StatusBadRequest},
		ValidationMinName:            {"El nombre de la empresa debe tener al menos 3 caracteres", http.StatusBadRequest},
		ValidationMaxName:            {"El nombre de la empresa debe tener como máximo 100 caracteres", http.StatusBadRequest},
		ValidationRequiredEmail:      {"El correo electrónico de la empresa es requerido", http.StatusBadRequest},
		ValidationEmailEmail:         {"El correo electrónico de la empresa debe ser válido", http.StatusBadRequest},
		ValidationRequiredPhone:      {"El teléfono de la empresa es requerido", http.StatusBadRequest},
		ValidationRequiredIdentifier: {"El identificador de la empresa es requerido", http.StatusBadRequest},
		ValidationRequiredLogo:       {"El logo de la empresa es requerido", http.StatusBadRequest},
	}

	return store
}

// GetMessage returns the message for the given key in the specified language
func (ms *MessageStore) GetMessage(lang, key string) (string, int) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	// Default to English if language not found
	if _, ok := ms.messages[lang]; !ok {
		lang = English
	}
	if msg, ok := ms.messages[lang][key]; ok {
		return msg.Text, msg.Status
	}
	return "", http.StatusInternalServerError
}

func (ms *MessageStore) AddLanguage(lang string, messages map[string]Message) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.messages[lang] = messages
}

func (ms *MessageStore) UpdateMessage(lang, key string, message Message) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, ok := ms.messages[lang]; !ok {
		ms.messages[lang] = make(map[string]Message)
	}

	ms.messages[lang][key] = message
}
