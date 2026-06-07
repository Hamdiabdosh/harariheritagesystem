package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleRegistrar  Role = "registrar"
	RoleSupervisor Role = "supervisor"
	RoleManager    Role = "manager"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleRegistrar, RoleSupervisor, RoleManager:
		return true
	default:
		return false
	}
}

type Language string

const (
	LanguageAm Language = "am"
	LanguageEn Language = "en"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	FullName     string     `json:"full_name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         Role       `json:"role"`
	Language     Language   `json:"language"`
	IsActive     bool       `json:"is_active"`
	DeletedAt    *time.Time `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (u *User) IsUsable() bool {
	return u.IsActive && u.DeletedAt == nil
}

type UserPublic struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	Role     Role      `json:"role"`
	Language Language  `json:"language"`
}

func (u *User) ToPublic() UserPublic {
	return UserPublic{
		ID:       u.ID,
		FullName: u.FullName,
		Email:    u.Email,
		Role:     u.Role,
		Language: u.Language,
	}
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
