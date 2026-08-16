package domain

import "time"

type Role string

const (
	RoleAdministrator Role = "administrator"
	RoleArchivist     Role = "archivist"
	RoleAuditor       Role = "auditor"
	RoleDepartment    Role = "department_manager"
	RoleBorrower      Role = "borrower"
)

type UserStatus string

const (
	UserActive   UserStatus = "active"
	UserDisabled UserStatus = "disabled"
)

type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	DisplayName  string     `json:"display_name"`
	AvatarURL    string     `json:"avatar_url,omitempty"`
	DepartmentID string     `json:"department_id"`
	Role         Role       `json:"role"`
	Status       UserStatus `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (u User) IsAdministrator() bool { return u.Role == RoleAdministrator }

type Department struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshToken struct {
	Hash      string
	UserID    string
	FamilyID  string
	ExpiresAt time.Time
	UsedAt    *time.Time
	RevokedAt *time.Time
}

func (t RefreshToken) CanRotate(now time.Time) error {
	if t.RevokedAt != nil {
		return ErrTokenRevoked
	}
	if t.UsedAt != nil {
		return ErrTokenReplayed
	}
	if !now.Before(t.ExpiresAt) {
		return ErrTokenRevoked
	}
	return nil
}
