package role

import "time"

type Role struct {
	ID          int64     `json:"id"`
	CompanyID   int64     `json:"company_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
