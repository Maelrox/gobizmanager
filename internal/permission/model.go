package permission

import "time"

type Permission struct {
	ID             int64         `json:"id"`
	RoleID         int64         `json:"role_id"`
	ModuleActionID int64         `json:"module_action_id"`
	ModuleAction   *ModuleAction `json:"module_action,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type ModuleAction struct {
	ID          int64     `json:"id"`
	ModuleID    int64     `json:"module_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
