package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/http/dto"
	httpmw "github.com/pafthang/bms/internal/http/middleware"
	serviceauth "github.com/pafthang/bms/internal/services/auth"
	"github.com/pafthang/bms/internal/services/workspaces"
	"github.com/pafthang/orm"
)

type workspacesHandler struct {
	svc *workspaces.Service
}

func RegisterWorkspaces(api *arc.Group, svc *workspaces.Service, auth *serviceauth.Service) {
	h := &workspacesHandler{svc: svc}
	secured := api.Group("", httpmw.RequireAuth(auth))
	ws := secured.Group("/workspaces").WithTags("Workspaces")
	members := secured.Group("/workspaces/{workspaceId}/users").WithTags("WorkspaceUsers")

	arc.HandleGroup(ws, http.MethodPost, "", "workspaces_create", h.create,
		arc.WithSecurity("BearerAuth"),
		arc.WithRequestExamples(map[string]any{
			"default": map[string]any{
				"name":        "Main Workspace",
				"description": "Personal bookmarks",
			},
		}),
		arc.WithResponseExamples(map[string]any{
			"default": map[string]any{
				"data": map[string]any{
					"workspace": map[string]any{
						"id":            10,
						"name":          "Main Workspace",
						"description":   "Personal bookmarks",
						"owner_user_id": 1,
						"created_at":    "2026-03-08T10:00:00Z",
						"updated_at":    "2026-03-08T10:00:00Z",
					},
					"membership": map[string]any{
						"id":           100,
						"workspace_id": 10,
						"user_id":      1,
						"role":         "admin",
						"email":        "user@example.com",
					},
				},
			},
		}),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(ws, http.MethodGet, "", "workspaces_list", h.list,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(ws, http.MethodGet, "/{workspaceId}", "workspaces_get", h.get,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized: {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusNotFound:     {Code: "workspace_not_found", Detail: "workspace not found"},
		}),
	)
	arc.HandleGroup(ws, http.MethodPatch, "/{workspaceId}", "workspaces_patch", h.patch,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:           {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusNotFound:            {Code: "workspace_not_found", Detail: "workspace not found"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(ws, http.MethodDelete, "/{workspaceId}", "workspaces_delete", h.delete,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized: {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:    {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusNotFound:     {Code: "workspace_not_found", Detail: "workspace not found"},
		}),
	)

	arc.HandleGroup(members, http.MethodGet, "", "workspace_users_list", h.listMembers,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized: {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:    {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusNotFound:     {Code: "workspace_not_found", Detail: "workspace not found"},
		}),
	)
	arc.HandleGroup(members, http.MethodPost, "", "workspace_users_add", h.addMember,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:           {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusNotFound:            {Code: "workspace_not_found", Detail: "workspace not found"},
			http.StatusConflict:            {Code: "workspace_membership_conflict", Detail: "membership already exists"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(members, http.MethodPatch, "/{userId}", "workspace_users_patch", h.patchMember,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:           {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusNotFound:            {Code: "workspace_not_found", Detail: "workspace not found"},
			http.StatusConflict:            {Code: "workspace_last_admin_violation", Detail: "cannot leave workspace without admin"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(members, http.MethodDelete, "/{userId}", "workspace_users_delete", h.deleteMember,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized: {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:    {Code: "workspace_forbidden", Detail: "insufficient role for workspace action"},
			http.StatusNotFound:     {Code: "workspace_not_found", Detail: "workspace not found"},
			http.StatusConflict:     {Code: "workspace_last_admin_violation", Detail: "cannot leave workspace without admin"},
		}),
	)
}

func (h *workspacesHandler) create(ctx context.Context, in *dto.WorkspaceCreateRequest) (*arc.Response[dto.DataEnvelope[dto.WorkspaceCreateResponse]], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	ws, m, err := h.svc.Create(ctx, id.UserID, workspaces.WorkspaceUpsertInput{Name: in.Name, Description: in.Description})
	if err != nil {
		return nil, mapWorkspaceError(err)
	}
	resp := dto.WorkspaceCreateResponse{
		Workspace: dto.ToWorkspaceResponse(ws),
		Membership: &dto.WorkspaceMembershipResponse{
			ID:          m.ID,
			WorkspaceID: m.WorkspaceID,
			UserID:      m.UserID,
			Role:        m.Role,
			Email:       id.Email,
		},
	}
	return arc.Created(dto.DataEnvelope[dto.WorkspaceCreateResponse]{Data: resp}), nil
}

func (h *workspacesHandler) list(ctx context.Context, in *dto.WorkspacesListRequest) (*arc.Response[dto.ListEnvelope[dto.WorkspaceResponse]], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	items, total, err := h.svc.List(ctx, workspaces.WorkspaceListInput{
		RequesterID: id.UserID,
		Superadmin:  id.IsSuperadmin,
		AllForAdmin: in.All,
		Search:      in.Search,
		Limit:       in.Limit,
		Offset:      in.Offset,
	})
	if err != nil {
		return nil, mapWorkspaceError(err)
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
	out := make([]dto.WorkspaceResponse, 0, len(items))
	for _, it := range items {
		out = append(out, *dto.ToWorkspaceResponse(&it))
	}
	return arc.OK(dto.ListEnvelope[dto.WorkspaceResponse]{Data: out, Meta: dto.BuildMeta(limit, offset, total)}), nil
}

func (h *workspacesHandler) get(ctx context.Context, in *dto.WorkspaceByIDRequest) (*arc.Response[dto.DataEnvelope[*dto.WorkspaceResponse]], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	ws, err := h.svc.Get(ctx, id.UserID, id.IsSuperadmin, in.WorkspaceID)
	if err != nil {
		return nil, mapWorkspaceError(err)
	}
	return arc.OK(dto.DataEnvelope[*dto.WorkspaceResponse]{Data: dto.ToWorkspaceResponse(ws)}), nil
}

func (h *workspacesHandler) patch(ctx context.Context, in *dto.WorkspacePatchByIDRequest) (*arc.Response[dto.DataEnvelope[*dto.WorkspaceResponse]], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	patch := workspaces.WorkspacePatchInput{
		NameSet:         in.Name.IsSet(),
		DescriptionSet:  in.Description.IsSet(),
		DescriptionNull: in.Description.IsNull(),
	}
	if v, ok := in.Name.Value(); ok {
		patch.Name = v
	}
	if v, ok := in.Description.Value(); ok {
		patch.Description = v
	}
	ws, err := h.svc.Patch(ctx, id.UserID, id.IsSuperadmin, in.WorkspaceID, patch)
	if err != nil {
		return nil, mapWorkspaceError(err)
	}
	return arc.OK(dto.DataEnvelope[*dto.WorkspaceResponse]{Data: dto.ToWorkspaceResponse(ws)}), nil
}

func (h *workspacesHandler) delete(ctx context.Context, in *dto.WorkspaceByIDRequest) (*arc.Response[struct{}], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	if err := h.svc.Delete(ctx, id.UserID, id.IsSuperadmin, in.WorkspaceID); err != nil {
		return nil, mapWorkspaceError(err)
	}
	return arc.NoContent(), nil
}

func (h *workspacesHandler) listMembers(ctx context.Context, in *dto.WorkspaceUsersListRequest) (*arc.Response[dto.ListEnvelope[dto.WorkspaceMembershipResponse]], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	items, total, err := h.svc.ListMembers(ctx, workspaces.MembershipListInput{
		RequesterID: id.UserID,
		Superadmin:  id.IsSuperadmin,
		WorkspaceID: in.WorkspaceID,
		Limit:       in.Limit,
		Offset:      in.Offset,
	})
	if err != nil {
		return nil, mapWorkspaceError(err)
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
	out := make([]dto.WorkspaceMembershipResponse, 0, len(items))
	for _, it := range items {
		out = append(out, dto.ToWorkspaceMembershipResponse(it))
	}
	return arc.OK(dto.ListEnvelope[dto.WorkspaceMembershipResponse]{Data: out, Meta: dto.BuildMeta(limit, offset, total)}), nil
}

func (h *workspacesHandler) addMember(ctx context.Context, in *dto.WorkspaceUserCreateRequest) (*arc.Response[dto.DataEnvelope[dto.WorkspaceMembershipResponse]], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	res, err := h.svc.AddMember(ctx, workspaces.MembershipCreateInput{
		RequesterID: id.UserID,
		Superadmin:  id.IsSuperadmin,
		WorkspaceID: in.WorkspaceID,
		Role:        in.Role,
		Email:       in.Email,
		UserID:      in.UserID,
	})
	if err != nil {
		return nil, mapWorkspaceError(err)
	}
	return arc.Created(dto.DataEnvelope[dto.WorkspaceMembershipResponse]{Data: dto.ToWorkspaceMembershipResponse(*res)}), nil
}

func (h *workspacesHandler) patchMember(ctx context.Context, in *dto.WorkspaceUserPatchRequest) (*arc.Response[dto.DataEnvelope[dto.WorkspaceMembershipResponse]], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	res, err := h.svc.UpdateMemberRole(ctx, workspaces.MembershipUpdateInput{
		RequesterID:  id.UserID,
		Superadmin:   id.IsSuperadmin,
		WorkspaceID:  in.WorkspaceID,
		TargetUserID: in.UserID,
		Role:         in.Role,
	})
	if err != nil {
		return nil, mapWorkspaceError(err)
	}
	return arc.OK(dto.DataEnvelope[dto.WorkspaceMembershipResponse]{Data: dto.ToWorkspaceMembershipResponse(*res)}), nil
}

func (h *workspacesHandler) deleteMember(ctx context.Context, in *dto.WorkspaceUserDeleteRequest) (*arc.Response[struct{}], error) {
	id, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, unauthorized()
	}
	if err := h.svc.DeleteMember(ctx, workspaces.MembershipDeleteInput{
		RequesterID:  id.UserID,
		Superadmin:   id.IsSuperadmin,
		WorkspaceID:  in.WorkspaceID,
		TargetUserID: in.UserID,
	}); err != nil {
		return nil, mapWorkspaceError(err)
	}
	return arc.NoContent(), nil
}

func mapWorkspaceError(err error) error {
	if err == nil {
		return nil
	}
	if orm.HasCode(err, orm.CodeNotFound) || errors.Is(err, workspaces.ErrWorkspaceNotVisible) {
		return &arc.APIError{Status: http.StatusNotFound, Code: "workspace_not_found", Message: "workspace not found"}
	}
	if errors.Is(err, workspaces.ErrForbidden) {
		return &arc.APIError{Status: http.StatusForbidden, Code: "workspace_forbidden", Message: "insufficient role for workspace action"}
	}
	if errors.Is(err, workspaces.ErrMembershipConflict) || orm.HasCode(err, orm.CodeConflict) {
		return &arc.APIError{Status: http.StatusConflict, Code: "workspace_membership_conflict", Message: "membership already exists"}
	}
	if errors.Is(err, workspaces.ErrLastAdminViolation) {
		return &arc.APIError{Status: http.StatusConflict, Code: "workspace_last_admin_violation", Message: "cannot leave workspace without admin"}
	}
	if errors.Is(err, workspaces.ErrInvalidRole) || errors.Is(err, workspaces.ErrValidation) {
		return &arc.APIError{Status: http.StatusUnprocessableEntity, Code: "validation_error", Message: "validation failed"}
	}
	if errors.Is(err, workspaces.ErrNotConfigured) {
		return &arc.APIError{Status: http.StatusInternalServerError, Code: "internal_error", Message: "service is not configured"}
	}
	if strings.Contains(strings.ToLower(err.Error()), "validation") {
		return &arc.APIError{Status: http.StatusUnprocessableEntity, Code: "validation_error", Message: "validation failed"}
	}
	return err
}

func unauthorized() error {
	return &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_unauthorized", Message: "missing auth context"}
}
