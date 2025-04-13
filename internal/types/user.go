package types

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
