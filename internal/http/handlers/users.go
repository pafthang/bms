package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/http/dto"
	httpmw "github.com/pafthang/bms/internal/http/middleware"
	serviceauth "github.com/pafthang/bms/internal/services/auth"
	settingssvc "github.com/pafthang/bms/internal/services/settings"
	userssvc "github.com/pafthang/bms/internal/services/users"
	"github.com/pafthang/orm"
)

type usersHandler struct {
	users    *userssvc.Service
	settings *settingssvc.Service
}

func RegisterUsers(api *arc.Group, users *userssvc.Service, settings *settingssvc.Service, auth *serviceauth.Service) {
	h := &usersHandler{users: users, settings: settings}
	secured := api.Group("", httpmw.RequireAuth(auth))
	usersGroup := secured.Group("/users").WithTags("Users")
	adminGroup := secured.Group("/admin").WithTags("Admin")

	arc.HandleGroup(usersGroup, http.MethodGet, "/me", "users_me_get", h.me,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized: {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusNotFound:     {Code: "not_found", Detail: "resource not found"},
		}),
	)
	arc.HandleGroup(usersGroup, http.MethodPatch, "/me", "users_me_patch", h.mePatch,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "email change is disabled in mvp"},
		}),
	)

	settingsGroup := secured.Group("/users/me/settings").WithTags("UserSettings")
	arc.HandleGroup(settingsGroup, http.MethodGet, "", "users_me_settings_get", h.settingsGet,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized: {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusNotFound:     {Code: "not_found", Detail: "resource not found"},
		}),
	)
	arc.HandleGroup(settingsGroup, http.MethodPut, "", "users_me_settings_put", h.settingsPut,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(settingsGroup, http.MethodPatch, "", "users_me_settings_patch", h.settingsPatch,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)

	arc.HandleGroup(adminGroup, http.MethodGet, "/users", "admin_users_list", h.adminUsersList,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized:        {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:           {Code: "workspace_forbidden", Detail: "superadmin required"},
			http.StatusUnprocessableEntity: {Code: "validation_error", Detail: "validation failed"},
		}),
	)
	arc.HandleGroup(adminGroup, http.MethodGet, "/users/{userId}", "admin_users_get", h.adminUsersGet,
		arc.WithSecurity("BearerAuth"),
		arc.WithProblemResponseSpec(map[int]arc.ProblemExampleSpec{
			http.StatusUnauthorized: {Code: "auth_unauthorized", Detail: "missing auth context"},
			http.StatusForbidden:    {Code: "workspace_forbidden", Detail: "superadmin required"},
			http.StatusNotFound:     {Code: "not_found", Detail: "resource not found"},
		}),
	)
}

func (h *usersHandler) me(ctx context.Context, in *struct{}) (*arc.Response[dto.DataEnvelope[*dto.UserMeResponse]], error) {
	identity, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_unauthorized", Message: "missing auth context"}
	}
	user, err := h.users.Me(ctx, identity.UserID)
	if err != nil {
		return nil, mapUserError(err)
	}
	return arc.OK(dto.DataEnvelope[*dto.UserMeResponse]{Data: dto.ToUserMeResponse(user)}), nil
}

func (h *usersHandler) settingsGet(ctx context.Context, in *struct{}) (*arc.Response[dto.DataEnvelope[*dto.UserSettingsResponse]], error) {
	identity, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_unauthorized", Message: "missing auth context"}
	}
	settings, err := h.settings.Get(ctx, identity.UserID)
	if err != nil {
		return nil, mapUserError(err)
	}
	return arc.OK(dto.DataEnvelope[*dto.UserSettingsResponse]{Data: dto.ToUserSettingsResponse(settings)}), nil
}

func (h *usersHandler) settingsPut(ctx context.Context, in *dto.PutUserSettingsRequest) (*arc.Response[dto.DataEnvelope[*dto.UserSettingsResponse]], error) {
	identity, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_unauthorized", Message: "missing auth context"}
	}
	settings, err := h.settings.Put(ctx, identity.UserID, settingssvc.PutInput{
		Theme:            in.Theme,
		Locale:           in.Locale,
		Timezone:         in.Timezone,
		BookmarksPerPage: in.BookmarksPerPage,
	})
	if err != nil {
		return nil, mapUserError(err)
	}
	return arc.OK(dto.DataEnvelope[*dto.UserSettingsResponse]{Data: dto.ToUserSettingsResponse(settings)}), nil
}

func (h *usersHandler) settingsPatch(ctx context.Context, in *dto.PatchUserSettingsRequest) (*arc.Response[dto.DataEnvelope[*dto.UserSettingsResponse]], error) {
	identity, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_unauthorized", Message: "missing auth context"}
	}
	settings, err := h.settings.Patch(ctx, identity.UserID, settingssvc.PatchInput{
		Theme:            in.Theme,
		Locale:           in.Locale,
		Timezone:         in.Timezone,
		BookmarksPerPage: in.BookmarksPerPage,
	})
	if err != nil {
		return nil, mapUserError(err)
	}
	return arc.OK(dto.DataEnvelope[*dto.UserSettingsResponse]{Data: dto.ToUserSettingsResponse(settings)}), nil
}

func (h *usersHandler) mePatch(ctx context.Context, in *dto.PatchUserMeRequest) (*arc.Response[dto.DataEnvelope[*dto.UserMeResponse]], error) {
	identity, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_unauthorized", Message: "missing auth context"}
	}
	patch := userssvc.PatchMeInput{EmailSet: in.Email.IsSet()}
	if v, ok := in.Email.Value(); ok {
		patch.Email = v
	}
	user, err := h.users.PatchMe(ctx, identity.UserID, patch)
	if err != nil {
		return nil, mapUserError(err)
	}
	return arc.OK(dto.DataEnvelope[*dto.UserMeResponse]{Data: dto.ToUserMeResponse(user)}), nil
}

func (h *usersHandler) adminUsersList(ctx context.Context, in *dto.AdminUsersListRequest) (*arc.Response[dto.ListEnvelope[dto.AdminUserResponse]], error) {
	identity, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_unauthorized", Message: "missing auth context"}
	}
	if !identity.IsSuperadmin {
		return nil, &arc.APIError{Status: http.StatusForbidden, Code: "workspace_forbidden", Message: "superadmin required"}
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
	items, total, err := h.users.AdminList(ctx, limit, offset, in.Search)
	if err != nil {
		return nil, mapUserError(err)
	}
	out := make([]dto.AdminUserResponse, 0, len(items))
	for i := range items {
		out = append(out, *dto.ToAdminUserResponse(&items[i]))
	}
	return arc.OK(dto.ListEnvelope[dto.AdminUserResponse]{Data: out, Meta: dto.BuildMeta(limit, offset, total)}), nil
}

func (h *usersHandler) adminUsersGet(ctx context.Context, in *dto.AdminUserByIDRequest) (*arc.Response[dto.DataEnvelope[*dto.AdminUserResponse]], error) {
	identity, ok := httpmw.IdentityFromContext(ctx)
	if !ok {
		return nil, &arc.APIError{Status: http.StatusUnauthorized, Code: "auth_unauthorized", Message: "missing auth context"}
	}
	if !identity.IsSuperadmin {
		return nil, &arc.APIError{Status: http.StatusForbidden, Code: "workspace_forbidden", Message: "superadmin required"}
	}
	user, err := h.users.AdminGet(ctx, in.UserID)
	if err != nil {
		return nil, mapUserError(err)
	}
	return arc.OK(dto.DataEnvelope[*dto.AdminUserResponse]{Data: dto.ToAdminUserResponse(user)}), nil
}

func mapUserError(err error) error {
	if err == nil {
		return nil
	}
	if orm.HasCode(err, orm.CodeNotFound) {
		return &arc.APIError{Status: http.StatusNotFound, Code: "not_found", Message: "resource not found"}
	}
	if errors.Is(err, userssvc.ErrEmailChangeForbidden) {
		return &arc.APIError{Status: http.StatusUnprocessableEntity, Code: "validation_error", Message: "email change is disabled in mvp"}
	}
	if errors.Is(err, userssvc.ErrNotConfigured) || errors.Is(err, settingssvc.ErrNotConfigured) {
		return &arc.APIError{Status: http.StatusInternalServerError, Code: "internal_error", Message: "service is not configured"}
	}
	return err
}
