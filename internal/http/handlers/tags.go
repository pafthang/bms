package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/http/dto"
	httpmw "github.com/pafthang/bms/internal/http/middleware"
	serviceauth "github.com/pafthang/bms/internal/services/auth"
	"github.com/pafthang/bms/internal/services/tags"
	"github.com/pafthang/orm"
)

type tagsHandler struct{ svc *tags.Service }

func RegisterTags(api *arc.Group, svc *tags.Service, auth *serviceauth.Service) {
	h := &tagsHandler{svc: svc}
	secured := api.Group("", httpmw.RequireAuth(auth))
	group := secured.Group("/workspaces/{workspaceId}/tags").WithTags("Tags")

	arc.HandleGroup(group, http.MethodGet, "", "tags_list", h.list,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:           {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(group, http.MethodPost, "", "tags_create", h.create,
		arc.WithSecurity("BearerAuth"),
		arc.WithRequestExamples(map[string]any{
			"default": map[string]any{
				"name":  "golang",
				"color": "#3b82f6",
			},
		}),
		arc.WithResponseExamples(map[string]any{
			"default": map[string]any{
				"data": map[string]any{
					"id":           7,
					"workspace_id": 10,
					"name":         "golang",
					"color":        "#3b82f6",
					"created_at":   "2026-03-08T10:20:00Z",
					"updated_at":   "2026-03-08T10:20:00Z",
				},
			},
		}),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:           {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusConflict:            {Code: "conflict_error", Detail: "resource conflict"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(group, http.MethodPatch, "/{tagId}", "tags_patch", h.patch,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:           {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusNotFound:            {Code: "tag_not_found", Detail: "tag not found"},
			http.StatusConflict:            {Code: "conflict_error", Detail: "resource conflict"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(group, http.MethodDelete, "/{tagId}", "tags_delete", h.delete,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized: {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:    {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusNotFound:     {Code: "tag_not_found", Detail: "tag not found"},
		}),
	)
}

func (h *tagsHandler) list(ctx context.Context, in *dto.TagsListRequest) (*arc.Response[dto.ListEnvelope[dto.TagResponse]], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	items, total, err := h.svc.List(ctx, tags.ListInput{RequesterID: id.UserID, Superadmin: id.IsSuperadmin, WorkspaceID: in.WorkspaceID, Limit: in.Limit, Offset: in.Offset})
	if err != nil {
		return nil, mapTagError(err)
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
	out := make([]dto.TagResponse, 0, len(items))
	for i := range items {
		out = append(out, *dto.ToTagResponse(&items[i]))
	}
	return arc.OK(dto.ListEnvelope[dto.TagResponse]{Data: out, Meta: dto.BuildMeta(limit, offset, total)}), nil
}

func (h *tagsHandler) create(ctx context.Context, in *dto.TagCreateRequest) (*arc.Response[dto.DataEnvelope[*dto.TagResponse]], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	if in.Color != nil && *in.Color != "" && !dto.IsValidHexColor(*in.Color) {
		return nil, &arc.APIError{Status: http.StatusUnprocessableEntity, Code: "validation_error", Message: "validation failed"}
	}
	item, err := h.svc.Create(ctx, tags.CreateInput{RequesterID: id.UserID, Superadmin: id.IsSuperadmin, WorkspaceID: in.WorkspaceID, Name: in.Name, Color: in.Color})
	if err != nil {
		return nil, mapTagError(err)
	}
	return arc.Created(dto.DataEnvelope[*dto.TagResponse]{Data: dto.ToTagResponse(item)}), nil
}

func (h *tagsHandler) patch(ctx context.Context, in *dto.TagPatchRequest) (*arc.Response[dto.DataEnvelope[*dto.TagResponse]], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	patch := tags.PatchInput{RequesterID: id.UserID, Superadmin: id.IsSuperadmin, WorkspaceID: in.WorkspaceID, TagID: in.TagID, NameSet: in.Name.IsSet(), ColorSet: in.Color.IsSet(), ColorNull: in.Color.IsNull()}
	if v, ok := in.Name.Value(); ok {
		patch.Name = v
	}
	if v, ok := in.Color.Value(); ok {
		if v != "" && !dto.IsValidHexColor(v) {
			return nil, &arc.APIError{Status: http.StatusUnprocessableEntity, Code: "validation_error", Message: "validation failed"}
		}
		patch.Color = v
	}
	item, err := h.svc.Patch(ctx, patch)
	if err != nil {
		return nil, mapTagError(err)
	}
	return arc.OK(dto.DataEnvelope[*dto.TagResponse]{Data: dto.ToTagResponse(item)}), nil
}

func (h *tagsHandler) delete(ctx context.Context, in *dto.TagDeleteRequest) (*arc.Response[struct{}], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	if err := h.svc.Delete(ctx, tags.DeleteInput{RequesterID: id.UserID, Superadmin: id.IsSuperadmin, WorkspaceID: in.WorkspaceID, TagID: in.TagID}); err != nil {
		return nil, mapTagError(err)
	}
	return arc.NoContent(), nil
}

func mapTagError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, tags.ErrForbidden) {
		return &arc.APIError{Status: http.StatusForbidden, Code: "workspace_forbidden", Message: "insufficient role for workspace action"}
	}
	if errors.Is(err, tags.ErrValidation) {
		return &arc.APIError{Status: http.StatusUnprocessableEntity, Code: "validation_error", Message: "validation failed"}
	}
	if errors.Is(err, tags.ErrNotConfigured) {
		return &arc.APIError{Status: http.StatusInternalServerError, Code: "internal_error", Message: "service is not configured"}
	}
	if orm.HasCode(err, orm.CodeConflict) {
		return &arc.APIError{Status: http.StatusConflict, Code: "conflict_error", Message: "resource conflict"}
	}
	if orm.HasCode(err, orm.CodeNotFound) {
		return &arc.APIError{Status: http.StatusNotFound, Code: "tag_not_found", Message: "tag not found"}
	}
	return err
}
