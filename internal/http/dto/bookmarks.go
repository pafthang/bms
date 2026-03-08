package dto

import (
	"time"

	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/services/bookmarks"
)

type BookmarksListRequest struct {
	WorkspaceID int64  `path:"workspaceId" validate:"required,min=1"`
	Limit       int64  `query:"limit" validate:"min=0,max=100"`
	Offset      int64  `query:"offset" validate:"min=0"`
	Search      string `query:"search"`
	Tag         string `query:"tag"`
	TagID       *int64 `query:"tag_id"`
	IsArchived  *bool  `query:"is_archived"`
	Sort        string `query:"sort"`
	Order       string `query:"order"`
}

type BookmarkCreateRequest struct {
	WorkspaceID int64   `path:"workspaceId" validate:"required,min=1"`
	Title       string  `json:"title" validate:"required,minlength=1,maxlength=300"`
	URL         string  `json:"url" validate:"required"`
	Description *string `json:"description,omitempty" validate:"omitempty,maxlength=5000"`
	IsArchived  *bool   `json:"is_archived,omitempty"`
	TagIDs      []int64 `json:"tag_ids,omitempty"`
}

type BookmarkByIDRequest struct {
	WorkspaceID int64 `path:"workspaceId" validate:"required,min=1"`
	BookmarkID  int64 `path:"bookmarkId" validate:"required,min=1"`
}

type BookmarkPatchRequest struct {
	WorkspaceID int64                 `path:"workspaceId" validate:"required,min=1"`
	BookmarkID  int64                 `path:"bookmarkId" validate:"required,min=1"`
	Title       arc.OptionalString    `json:"title,omitempty" validate:"omitempty,minlength=1,maxlength=300"`
	URL         arc.OptionalString    `json:"url,omitempty"`
	Description arc.OptionalString    `json:"description,omitempty" validate:"omitempty,maxlength=5000"`
	IsArchived  arc.OptionalBool      `json:"is_archived,omitempty"`
	TagIDs      arc.Optional[[]int64] `json:"tag_ids,omitempty"`
}

type BookmarkResponse struct {
	ID              int64   `json:"id"`
	WorkspaceID     int64   `json:"workspace_id"`
	CreatedByUserID int64   `json:"created_by_user_id"`
	UpdatedByUserID *int64  `json:"updated_by_user_id,omitempty"`
	Title           string  `json:"title"`
	URL             string  `json:"url"`
	Description     *string `json:"description,omitempty"`
	IsArchived      bool    `json:"is_archived"`
	TagIDs          []int64 `json:"tag_ids,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

func ToBookmarkResponse(v *bookmarks.BookmarkView) *BookmarkResponse {
	b := v.Model
	return &BookmarkResponse{
		ID:              b.ID,
		WorkspaceID:     b.WorkspaceID,
		CreatedByUserID: b.CreatedByUserID,
		UpdatedByUserID: b.UpdatedByUserID,
		Title:           b.Title,
		URL:             b.URL,
		Description:     b.Description,
		IsArchived:      b.IsArchived,
		TagIDs:          append([]int64{}, v.TagIDs...),
		CreatedAt:       b.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:       b.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
