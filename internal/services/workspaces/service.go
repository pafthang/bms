package workspaces

import (
	"context"
	"errors"
	"strings"

	"github.com/pafthang/bms/internal/domain/models"
	"github.com/pafthang/dbx"
	"github.com/pafthang/orm"
)

const (
	RoleViewer = "viewer"
	RoleEditor = "editor"
	RoleAdmin  = "admin"
)

var (
	ErrNotConfigured       = errors.New("database is not configured")
	ErrForbidden           = errors.New("forbidden")
	ErrInvalidRole         = errors.New("invalid workspace role")
	ErrMembershipConflict  = errors.New("workspace membership already exists")
	ErrLastAdminViolation  = errors.New("cannot remove or demote last admin")
	ErrWorkspaceNotVisible = errors.New("workspace not visible")
	ErrValidation          = errors.New("validation error")
)

type Service struct {
	db *dbx.DB
}

type WorkspaceListInput struct {
	RequesterID int64
	Superadmin  bool
	AllForAdmin bool
	Search      string
	Limit       int64
	Offset      int64
}

type WorkspaceUpsertInput struct {
	Name        string
	Description *string
}

type WorkspacePatchInput struct {
	NameSet         bool
	Name            string
	DescriptionSet  bool
	DescriptionNull bool
	Description     string
}

type MembershipListInput struct {
	RequesterID int64
	Superadmin  bool
	WorkspaceID int64
	Limit       int64
	Offset      int64
}

type MembershipCreateInput struct {
	RequesterID int64
	Superadmin  bool
	WorkspaceID int64
	Role        string
	Email       *string
	UserID      *int64
}

type MembershipUpdateInput struct {
	RequesterID  int64
	Superadmin   bool
	WorkspaceID  int64
	TargetUserID int64
	Role         string
}

type MembershipDeleteInput struct {
	RequesterID  int64
	Superadmin   bool
	WorkspaceID  int64
	TargetUserID int64
}

type MembershipView struct {
	ID          int64
	WorkspaceID int64
	UserID      int64
	Email       string
	Role        string
}

func NewService(db *dbx.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(ctx context.Context, requesterID int64, in WorkspaceUpsertInput) (*models.Workspace, *models.WorkspaceUser, error) {
	if s.db == nil {
		return nil, nil, ErrNotConfigured
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, nil, ErrValidation
	}
	ws := &models.Workspace{Name: name, Description: in.Description, OwnerUserID: requesterID}
	membership := &models.WorkspaceUser{UserID: requesterID, Role: RoleAdmin}
	if err := orm.WithTx(ctx, s.db, func(tx *orm.Tx) error {
		if err := orm.Insert(ctx, tx, ws); err != nil {
			return err
		}
		membership.WorkspaceID = ws.ID
		if err := orm.Insert(ctx, tx, membership); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, nil, err
	}
	return ws, membership, nil
}

func (s *Service) List(ctx context.Context, in WorkspaceListInput) ([]models.Workspace, int64, error) {
	if s.db == nil {
		return nil, 0, ErrNotConfigured
	}
	limit, offset := normalizePaging(in.Limit, in.Offset)
	search := strings.TrimSpace(in.Search)

	q := orm.Query[models.Workspace](s.db)
	if in.Superadmin && in.AllForAdmin {
		// unrestricted list
	} else {
		memberships, err := orm.Query[models.WorkspaceUser](s.db).
			WhereEq("user_id", in.RequesterID).
			All(ctx)
		if err != nil {
			return nil, 0, err
		}
		if len(memberships) == 0 {
			return []models.Workspace{}, 0, nil
		}
		ids := make([]int64, 0, len(memberships))
		for _, m := range memberships {
			ids = append(ids, m.WorkspaceID)
		}
		q = q.WhereIn("id", ids)
	}
	if search != "" {
		q = q.WhereLike("name", "%"+search+"%")
	}
	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	items, err := q.OrderByDesc("created_at").OrderByDesc("id").Limit(limit).Offset(offset).All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *Service) Get(ctx context.Context, requesterID int64, superadmin bool, workspaceID int64) (*models.Workspace, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	ws, err := orm.ByPK[models.Workspace](ctx, s.db, workspaceID)
	if err != nil {
		return nil, err
	}
	if superadmin {
		return ws, nil
	}
	ok, err := s.isMember(ctx, workspaceID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrWorkspaceNotVisible
	}
	return ws, nil
}

func (s *Service) Patch(ctx context.Context, requesterID int64, superadmin bool, workspaceID int64, in WorkspacePatchInput) (*models.Workspace, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	if err := s.requireWorkspaceRole(ctx, requesterID, superadmin, workspaceID, RoleAdmin); err != nil {
		return nil, err
	}
	ws, err := orm.ByPK[models.Workspace](ctx, s.db, workspaceID)
	if err != nil {
		return nil, err
	}
	if in.NameSet {
		name := strings.TrimSpace(in.Name)
		if name == "" {
			return nil, ErrValidation
		}
		ws.Name = name
	}
	if in.DescriptionSet {
		if in.DescriptionNull {
			ws.Description = nil
		} else {
			d := in.Description
			ws.Description = &d
		}
	}
	if err := orm.Update(ctx, s.db, ws); err != nil {
		return nil, err
	}
	return ws, nil
}

func (s *Service) Delete(ctx context.Context, requesterID int64, superadmin bool, workspaceID int64) error {
	if s.db == nil {
		return ErrNotConfigured
	}
	if err := s.requireWorkspaceRole(ctx, requesterID, superadmin, workspaceID, RoleAdmin); err != nil {
		return err
	}
	return orm.DeleteByPK[models.Workspace](ctx, s.db, workspaceID)
}

func (s *Service) ListMembers(ctx context.Context, in MembershipListInput) ([]MembershipView, int64, error) {
	if s.db == nil {
		return nil, 0, ErrNotConfigured
	}
	if err := s.requireWorkspaceMember(ctx, in.RequesterID, in.Superadmin, in.WorkspaceID); err != nil {
		return nil, 0, err
	}
	limit, offset := normalizePaging(in.Limit, in.Offset)
	q := orm.Query[models.WorkspaceUser](s.db).WhereEq("workspace_id", in.WorkspaceID)
	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	items, err := q.OrderBy("created_at").OrderBy("id").Limit(limit).Offset(offset).All(ctx)
	if err != nil {
		return nil, 0, err
	}
	out := make([]MembershipView, 0, len(items))
	for _, m := range items {
		email := ""
		u, err := orm.ByPK[models.User](ctx, s.db, m.UserID)
		if err == nil {
			email = u.Email
		}
		out = append(out, MembershipView{ID: m.ID, WorkspaceID: m.WorkspaceID, UserID: m.UserID, Email: email, Role: m.Role})
	}
	return out, total, nil
}

func (s *Service) AddMember(ctx context.Context, in MembershipCreateInput) (*MembershipView, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	if err := s.requireWorkspaceRole(ctx, in.RequesterID, in.Superadmin, in.WorkspaceID, RoleAdmin); err != nil {
		return nil, err
	}
	role, err := normalizeRole(in.Role)
	if err != nil {
		return nil, err
	}
	targetID, targetEmail, err := s.resolveTargetUser(ctx, in.UserID, in.Email)
	if err != nil {
		return nil, err
	}
	if _, err := orm.Query[models.WorkspaceUser](s.db).
		WhereEq("workspace_id", in.WorkspaceID).
		WhereEq("user_id", targetID).
		One(ctx); err == nil {
		return nil, ErrMembershipConflict
	} else if !orm.HasCode(err, orm.CodeNotFound) {
		return nil, err
	}
	m := &models.WorkspaceUser{WorkspaceID: in.WorkspaceID, UserID: targetID, Role: role}
	if err := orm.Insert(ctx, s.db, m); err != nil {
		return nil, err
	}
	return &MembershipView{ID: m.ID, WorkspaceID: m.WorkspaceID, UserID: m.UserID, Email: targetEmail, Role: m.Role}, nil
}

func (s *Service) UpdateMemberRole(ctx context.Context, in MembershipUpdateInput) (*MembershipView, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	if err := s.requireWorkspaceRole(ctx, in.RequesterID, in.Superadmin, in.WorkspaceID, RoleAdmin); err != nil {
		return nil, err
	}
	role, err := normalizeRole(in.Role)
	if err != nil {
		return nil, err
	}
	member, err := orm.Query[models.WorkspaceUser](s.db).
		WhereEq("workspace_id", in.WorkspaceID).
		WhereEq("user_id", in.TargetUserID).
		One(ctx)
	if err != nil {
		return nil, err
	}
	if member.Role == RoleAdmin && role != RoleAdmin {
		count, err := orm.Query[models.WorkspaceUser](s.db).
			WhereEq("workspace_id", in.WorkspaceID).
			WhereEq("role", RoleAdmin).
			WhereNotEq("user_id", in.TargetUserID).
			Count(ctx)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, ErrLastAdminViolation
		}
	}
	member.Role = role
	if err := orm.Update(ctx, s.db, member); err != nil {
		return nil, err
	}
	email := ""
	if u, err := orm.ByPK[models.User](ctx, s.db, member.UserID); err == nil {
		email = u.Email
	}
	return &MembershipView{ID: member.ID, WorkspaceID: member.WorkspaceID, UserID: member.UserID, Email: email, Role: member.Role}, nil
}

func (s *Service) DeleteMember(ctx context.Context, in MembershipDeleteInput) error {
	if s.db == nil {
		return ErrNotConfigured
	}
	if err := s.requireWorkspaceRole(ctx, in.RequesterID, in.Superadmin, in.WorkspaceID, RoleAdmin); err != nil {
		return err
	}
	member, err := orm.Query[models.WorkspaceUser](s.db).
		WhereEq("workspace_id", in.WorkspaceID).
		WhereEq("user_id", in.TargetUserID).
		One(ctx)
	if err != nil {
		return err
	}
	if member.Role == RoleAdmin {
		count, err := orm.Query[models.WorkspaceUser](s.db).
			WhereEq("workspace_id", in.WorkspaceID).
			WhereEq("role", RoleAdmin).
			WhereNotEq("user_id", in.TargetUserID).
			Count(ctx)
		if err != nil {
			return err
		}
		if count == 0 {
			return ErrLastAdminViolation
		}
	}
	return orm.Delete(ctx, s.db, member)
}

func (s *Service) resolveTargetUser(ctx context.Context, userID *int64, email *string) (int64, string, error) {
	if userID != nil {
		u, err := orm.ByPK[models.User](ctx, s.db, *userID)
		if err != nil {
			return 0, "", err
		}
		return u.ID, u.Email, nil
	}
	if email != nil {
		norm := strings.ToLower(strings.TrimSpace(*email))
		if norm != "" {
			u, err := orm.Query[models.User](s.db).WhereEq("email", norm).One(ctx)
			if err != nil {
				return 0, "", err
			}
			return u.ID, u.Email, nil
		}
	}
	return 0, "", ErrValidation
}

func (s *Service) requireWorkspaceMember(ctx context.Context, requesterID int64, superadmin bool, workspaceID int64) error {
	if superadmin {
		return nil
	}
	ok, err := s.isMember(ctx, workspaceID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrForbidden
	}
	return nil
}

func (s *Service) requireWorkspaceRole(ctx context.Context, requesterID int64, superadmin bool, workspaceID int64, minRole string) error {
	if superadmin {
		return nil
	}
	member, err := orm.Query[models.WorkspaceUser](s.db).
		WhereEq("workspace_id", workspaceID).
		WhereEq("user_id", requesterID).
		One(ctx)
	if err != nil {
		if orm.HasCode(err, orm.CodeNotFound) {
			return ErrForbidden
		}
		return err
	}
	if roleRank(member.Role) < roleRank(minRole) {
		return ErrForbidden
	}
	return nil
}

func (s *Service) isMember(ctx context.Context, workspaceID, userID int64) (bool, error) {
	_, err := orm.Query[models.WorkspaceUser](s.db).
		WhereEq("workspace_id", workspaceID).
		WhereEq("user_id", userID).
		One(ctx)
	if err == nil {
		return true, nil
	}
	if orm.HasCode(err, orm.CodeNotFound) {
		return false, nil
	}
	return false, err
}

func normalizeRole(v string) (string, error) {
	r := strings.ToLower(strings.TrimSpace(v))
	switch r {
	case RoleViewer, RoleEditor, RoleAdmin:
		return r, nil
	default:
		return "", ErrInvalidRole
	}
}

func roleRank(role string) int {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case RoleAdmin:
		return 3
	case RoleEditor:
		return 2
	case RoleViewer:
		return 1
	default:
		return 0
	}
}

func normalizePaging(limit, offset int64) (int64, int64) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
