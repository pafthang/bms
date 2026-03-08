package dto

import (
	"time"

	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/domain/models"
	"github.com/pafthang/bms/internal/services/workspaces"
)

type WorkspaceResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	OwnerUserID int64   `json:"owner_user_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type WorkspaceMembershipResponse struct {
	ID          int64  `json:"id"`
	WorkspaceID int64  `json:"workspace_id"`
	UserID      int64  `json:"user_id"`
	Role        string `json:"role"`
	Email       string `json:"email,omitempty"`
}

type WorkspaceCreateRequest struct {
	Name        string  `json:"name" validate:"required,minlength=1,maxlength=120"`
	Description *string `json:"description,omitempty" validate:"omitempty,maxlength=2000"`
}

type WorkspacePatchRequest struct {
	Name        arc.OptionalString `json:"name,omitempty" validate:"omitempty,minlength=1,maxlength=120"`
	Description arc.OptionalString `json:"description,omitempty" validate:"omitempty,maxlength=2000"`
}

type WorkspacePatchByIDRequest struct {
	WorkspaceID int64              `path:"workspaceId" validate:"required,min=1"`
	Name        arc.OptionalString `json:"name,omitempty" validate:"omitempty,minlength=1,maxlength=120"`
	Description arc.OptionalString `json:"description,omitempty" validate:"omitempty,maxlength=2000"`
}

type WorkspacesListRequest struct {
	Limit  int64  `query:"limit" validate:"min=0,max=100"`
	Offset int64  `query:"offset" validate:"min=0"`
	Search string `query:"search"`
	All    bool   `query:"all"`
}

type WorkspaceByIDRequest struct {
	WorkspaceID int64 `path:"workspaceId" validate:"required,min=1"`
}

type WorkspaceUsersListRequest struct {
	WorkspaceID int64 `path:"workspaceId" validate:"required,min=1"`
	Limit       int64 `query:"limit" validate:"min=0,max=100"`
	Offset      int64 `query:"offset" validate:"min=0"`
}

type WorkspaceUserCreateRequest struct {
	WorkspaceID int64   `path:"workspaceId" validate:"required,min=1"`
	Email       *string `json:"email,omitempty"`
	UserID      *int64  `json:"user_id,omitempty"`
	Role        string  `json:"role" validate:"required,enum=viewer|editor|admin"`
}

type WorkspaceUserPatchRequest struct {
	WorkspaceID int64  `path:"workspaceId" validate:"required,min=1"`
	UserID      int64  `path:"userId" validate:"required,min=1"`
	Role        string `json:"role" validate:"required,enum=viewer|editor|admin"`
}

type WorkspaceUserDeleteRequest struct {
	WorkspaceID int64 `path:"workspaceId" validate:"required,min=1"`
	UserID      int64 `path:"userId" validate:"required,min=1"`
}

type WorkspaceCreateResponse struct {
	Workspace  *WorkspaceResponse           `json:"workspace"`
	Membership *WorkspaceMembershipResponse `json:"membership"`
}

func ToWorkspaceResponse(m *models.Workspace) *WorkspaceResponse {
	return &WorkspaceResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		OwnerUserID: m.OwnerUserID,
		CreatedAt:   m.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func ToWorkspaceMembershipResponse(m workspaces.MembershipView) WorkspaceMembershipResponse {
	return WorkspaceMembershipResponse{
		ID:          m.ID,
		WorkspaceID: m.WorkspaceID,
		UserID:      m.UserID,
		Role:        m.Role,
		Email:       m.Email,
	}
}

func BuildMeta(limit, offset, total int64) Meta {
	hasNext := offset+limit < total
	if limit == 0 {
		hasNext = false
	}
	return Meta{Limit: limit, Offset: offset, Total: total, HasNext: hasNext}
}
