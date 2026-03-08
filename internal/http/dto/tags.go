package dto

import (
	"regexp"
	"time"

	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/domain/models"
)

var hexColorRe = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$`)

type TagsListRequest struct {
	WorkspaceID int64 `path:"workspaceId" validate:"required,min=1"`
	Limit       int64 `query:"limit" validate:"min=0,max=100"`
	Offset      int64 `query:"offset" validate:"min=0"`
}

type TagCreateRequest struct {
	WorkspaceID int64   `path:"workspaceId" validate:"required,min=1"`
	Name        string  `json:"name" validate:"required,minlength=1,maxlength=50"`
	Color       *string `json:"color,omitempty"`
}

type TagPatchRequest struct {
	WorkspaceID int64              `path:"workspaceId" validate:"required,min=1"`
	TagID       int64              `path:"tagId" validate:"required,min=1"`
	Name        arc.OptionalString `json:"name,omitempty" validate:"omitempty,minlength=1,maxlength=50"`
	Color       arc.OptionalString `json:"color,omitempty"`
}

type TagDeleteRequest struct {
	WorkspaceID int64 `path:"workspaceId" validate:"required,min=1"`
	TagID       int64 `path:"tagId" validate:"required,min=1"`
}

type TagResponse struct {
	ID          int64   `json:"id"`
	WorkspaceID int64   `json:"workspace_id"`
	Name        string  `json:"name"`
	Color       *string `json:"color,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func ToTagResponse(t *models.Tag) *TagResponse {
	return &TagResponse{
		ID:          t.ID,
		WorkspaceID: t.WorkspaceID,
		Name:        t.Name,
		Color:       t.Color,
		CreatedAt:   t.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func IsValidHexColor(v string) bool {
	return hexColorRe.MatchString(v)
}
