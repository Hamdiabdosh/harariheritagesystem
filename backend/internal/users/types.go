package users

import (
	"time"

	"github.com/qirs-mezgeb/api/internal/models"
)

type ListFilters struct {
	Page     int
	Limit    int
	Role     models.Role
	IsActive *bool
}

type UserItem struct {
	models.UserPublic
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toUserItem(user *models.User) UserItem {
	return UserItem{
		UserPublic: user.ToPublic(),
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

type PaginatedUsers struct {
	Items      []UserItem `json:"items"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalPages int        `json:"total_pages"`
}

type CreateInput struct {
	FullName string
	Email    string
	Password string
	Role     models.Role
	Language models.Language
}

type UpdateInput struct {
	FullName *string
	Role     *models.Role
	Language *models.Language
	IsActive *bool
}
