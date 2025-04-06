package company

type CreateCompanyRequest struct {
	Name       string `json:"name" validate:"required,min=3,max=100" msg:"company.name_required"`
	Email      string `json:"email" validate:"required,email" msg:"company.email_required"`
	Phone      string `json:"phone" validate:"required" msg:"company.phone_required"`
	Address    string `json:"address" validate:"required" msg:"company.address_required"`
	Logo       string `json:"logo"`
	Identifier string `json:"identifier" validate:"required" msg:"company.identifier_required"`
}

type UpdateCompanyRequest struct {
	Name       string `json:"name" validate:"required,min=3,max=100" msg:"company.name_required"`
	Email      string `json:"email" validate:"required,email" msg:"company.email_required"`
	Phone      string `json:"phone" validate:"required" msg:"company.phone_required"`
	Identifier string `json:"identifier" validate:"required" msg:"company.identifier_required"`
}

type UpdateCompanyLogoRequest struct {
	Logo string `json:"logo" validate:"required" msg:"company.logo_required"`
}
