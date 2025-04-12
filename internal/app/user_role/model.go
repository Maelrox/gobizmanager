package user_role

import "time"

type UserRole struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id" gorm:"index"`
	RoleID    int64     `json:"role_id" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
