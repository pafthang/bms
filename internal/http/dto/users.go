package dto

import (
	"time"

	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/domain/models"
)

type UserMeResponse struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	IsSuperadmin bool   `json:"is_superadmin"`
	Status       string `json:"status"`
}

func ToUserMeResponse(m *models.User) *UserMeResponse {
	return &UserMeResponse{
		ID:           m.ID,
		Email:        m.Email,
		IsSuperadmin: m.IsSuperadmin,
		Status:       m.Status,
	}
}

type UserSettingsResponse struct {
	Theme            *string `json:"theme,omitempty"`
	Locale           *string `json:"locale,omitempty"`
	Timezone         *string `json:"timezone,omitempty"`
	BookmarksPerPage *int    `json:"bookmarks_per_page,omitempty"`
}

func ToUserSettingsResponse(s *models.UserSettings) *UserSettingsResponse {
	return &UserSettingsResponse{
		Theme:            s.Theme,
		Locale:           s.Locale,
		Timezone:         s.Timezone,
		BookmarksPerPage: s.BookmarksPerPage,
	}
}

type PutUserSettingsRequest struct {
	Theme            *string `json:"theme,omitempty"`
	Locale           *string `json:"locale,omitempty" validate:"omitempty,maxlength=16"`
	Timezone         *string `json:"timezone,omitempty" validate:"omitempty,maxlength=64"`
	BookmarksPerPage *int    `json:"bookmarks_per_page,omitempty" validate:"omitempty,min=10,max=200"`
}

type PatchUserSettingsRequest struct {
	Theme            arc.OptionalString `json:"theme,omitempty"`
	Locale           arc.OptionalString `json:"locale,omitempty" validate:"omitempty,maxlength=16"`
	Timezone         arc.OptionalString `json:"timezone,omitempty" validate:"omitempty,maxlength=64"`
	BookmarksPerPage arc.OptionalInt    `json:"bookmarks_per_page,omitempty" validate:"omitempty,min=10,max=200"`
}

type PatchUserMeRequest struct {
	Email arc.OptionalString `json:"email,omitempty" validate:"omitempty,format=email"`
}

type AdminUsersListRequest struct {
	Limit  int64  `query:"limit" validate:"min=0,max=100"`
	Offset int64  `query:"offset" validate:"min=0"`
	Search string `query:"search"`
}

type AdminUserByIDRequest struct {
	UserID int64 `path:"userId" validate:"required,min=1"`
}

type AdminUserResponse struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	IsSuperadmin bool   `json:"is_superadmin"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func ToAdminUserResponse(m *models.User) *AdminUserResponse {
	return &AdminUserResponse{
		ID:           m.ID,
		Email:        m.Email,
		IsSuperadmin: m.IsSuperadmin,
		Status:       m.Status,
		CreatedAt:    m.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:    m.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
