package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/http/dto"
	httpmw "github.com/pafthang/bms/internal/http/middleware"
	serviceauth "github.com/pafthang/bms/internal/services/auth"
	"github.com/pafthang/bms/internal/services/bookmarks"
	"github.com/pafthang/orm"
)

type bookmarksHandler struct{ svc *bookmarks.Service }

func RegisterBookmarks(api *arc.Group, svc *bookmarks.Service, auth *serviceauth.Service) {
	h := &bookmarksHandler{svc: svc}
	secured := api.Group("", httpmw.RequireAuth(auth))
	group := secured.Group("/workspaces/{workspaceId}/bookmarks").WithTags("Bookmarks")

	arc.HandleGroup(group, http.MethodGet, "", "bookmarks_list", h.list,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:           {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(group, http.MethodPost, "", "bookmarks_create", h.create,
		arc.WithSecurity("BearerAuth"),
		arc.WithRequestExamples(map[string]any{
			"default": map[string]any{
				"title":       "Arc repository",
				"url":         "https://github.com/pafthang/arc",
				"description": "Core HTTP runtime",
				"is_archived": false,
				"tag_ids":     []int64{1, 2},
			},
		}),
		arc.WithResponseExamples(map[string]any{
			"default": map[string]any{
				"data": map[string]any{
					"id":                 501,
					"workspace_id":       10,
					"created_by_user_id": 1,
					"title":              "Arc repository",
					"url":                "https://github.com/pafthang/arc",
					"description":        "Core HTTP runtime",
					"is_archived":        false,
					"tag_ids":            []int64{1, 2},
					"created_at":         "2026-03-08T10:10:00Z",
					"updated_at":         "2026-03-08T10:10:00Z",
				},
			},
		}),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:           {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(group, http.MethodGet, "/{bookmarkId}", "bookmarks_get", h.get,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized: {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:    {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusNotFound:     {Code: "bookmark_not_found", Detail: "bookmark not found"},
		}),
	)
	arc.HandleGroup(group, http.MethodPatch, "/{bookmarkId}", "bookmarks_patch", h.patch,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:           {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusNotFound:            {Code: "bookmark_not_found", Detail: "bookmark not found"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(group, http.MethodDelete, "/{bookmarkId}", "bookmarks_delete", h.delete,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized: {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:    {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusNotFound:     {Code: "bookmark_not_found", Detail: "bookmark not found"},
		}),
	)
}

func (h *bookmarksHandler) list(ctx context.Context, in *dto.BookmarksListRequest) (*arc.Response[dto.ListEnvelope[dto.BookmarkResponse]], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	items, total, err := h.svc.List(ctx, bookmarks.ListInput{
		RequesterID: id.UserID,
		Superadmin:  id.IsSuperadmin,
		WorkspaceID: in.WorkspaceID,
		Limit:       in.Limit,
		Offset:      in.Offset,
		Search:      in.Search,
		Tag:         in.Tag,
		TagID:       in.TagID,
		IsArchived:  in.IsArchived,
		Sort:        in.Sort,
		Order:       in.Order,
	})
	if err != nil {
		return nil, mapBookmarkError(err)
	}
	limit := in.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := in.Offset
	if offset < 0 {
		offset = 0
	}
	out := make([]dto.BookmarkResponse, 0, len(items))
	for i := range items {
		out = append(out, *dto.ToBookmarkResponse(&items[i]))
	}
	return arc.OK(dto.ListEnvelope[dto.BookmarkResponse]{Data: out, Meta: dto.BuildMeta(limit, offset, total)}), nil
}

func (h *bookmarksHandler) create(ctx context.Context, in *dto.BookmarkCreateRequest) (*arc.Response[dto.DataEnvelope[*dto.BookmarkResponse]], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	item, err := h.svc.Create(ctx, bookmarks.CreateInput{
		RequesterID: id.UserID,
		Superadmin:  id.IsSuperadmin,
		WorkspaceID: in.WorkspaceID,
		Title:       in.Title,
		URL:         in.URL,
		Description: in.Description,
		IsArchived:  in.IsArchived,
		TagIDs:      in.TagIDs,
	})
	if err != nil {
		return nil, mapBookmarkError(err)
	}
	return arc.Created(dto.DataEnvelope[*dto.BookmarkResponse]{Data: dto.ToBookmarkResponse(item)}), nil
}

func (h *bookmarksHandler) get(ctx context.Context, in *dto.BookmarkByIDRequest) (*arc.Response[dto.DataEnvelope[*dto.BookmarkResponse]], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	item, err := h.svc.Get(ctx, id.UserID, id.IsSuperadmin, in.WorkspaceID, in.BookmarkID)
	if err != nil {
		return nil, mapBookmarkError(err)
	}
	return arc.OK(dto.DataEnvelope[*dto.BookmarkResponse]{Data: dto.ToBookmarkResponse(item)}), nil
}

func (h *bookmarksHandler) patch(ctx context.Context, in *dto.BookmarkPatchRequest) (*arc.Response[dto.DataEnvelope[*dto.BookmarkResponse]], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	patch := bookmarks.PatchInput{
		RequesterID:     id.UserID,
		Superadmin:      id.IsSuperadmin,
		WorkspaceID:     in.WorkspaceID,
		BookmarkID:      in.BookmarkID,
		TitleSet:        in.Title.IsSet(),
		URLSet:          in.URL.IsSet(),
		DescriptionSet:  in.Description.IsSet(),
		DescriptionNull: in.Description.IsNull(),
		IsArchivedSet:   in.IsArchived.IsSet(),
		TagIDsSet:       in.TagIDs.IsSet(),
		UpdatedByUserID: id.UserID,
	}
	if v, ok := in.Title.Value(); ok {
		patch.Title = v
	}
	if v, ok := in.URL.Value(); ok {
		patch.URL = v
	}
	if v, ok := in.Description.Value(); ok {
		patch.Description = v
	}
	if v, ok := in.IsArchived.Value(); ok {
		patch.IsArchived = v
	}
	if v, ok := in.TagIDs.Value(); ok {
		patch.TagIDs = v
	}
	item, err := h.svc.Patch(ctx, patch)
	if err != nil {
		return nil, mapBookmarkError(err)
	}
	return arc.OK(dto.DataEnvelope[*dto.BookmarkResponse]{Data: dto.ToBookmarkResponse(item)}), nil
}

func (h *bookmarksHandler) delete(ctx context.Context, in *dto.BookmarkByIDRequest) (*arc.Response[struct{}], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	if err := h.svc.Delete(ctx, bookmarks.DeleteInput{RequesterID: id.UserID, Superadmin: id.IsSuperadmin, WorkspaceID: in.WorkspaceID, BookmarkID: in.BookmarkID}); err != nil {
		return nil, mapBookmarkError(err)
	}
	return arc.NoContent(), nil
}

func mapBookmarkError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, bookmarks.ErrForbidden) {
		return &arc.APIError{Status: http.StatusForbidden, Code: "workspace_forbidden", Message: "insufficient role for workspace action"}
	}
	if errors.Is(err, bookmarks.ErrValidation) {
		return &arc.APIError{Status: http.StatusUnprocessableEntity, Code: "validation_error", Message: "validation failed"}
	}
	if errors.Is(err, bookmarks.ErrNotConfigured) {
		return &arc.APIError{Status: http.StatusInternalServerError, Code: "internal_error", Message: "service is not configured"}
	}
	if orm.HasCode(err, orm.CodeNotFound) {
		return &arc.APIError{Status: http.StatusNotFound, Code: "bookmark_not_found", Message: "bookmark not found"}
	}
	if orm.HasCode(err, orm.CodeConflict) {
		return &arc.APIError{Status: http.StatusConflict, Code: "conflict_error", Message: "resource conflict"}
	}
	return err
}
