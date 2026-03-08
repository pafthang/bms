package tags

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
	ErrNotConfigured = errors.New("database is not configured")
	ErrForbidden     = errors.New("forbidden")
	ErrValidation    = errors.New("validation error")
)

type Service struct {
	db *dbx.DB
}

type ListInput struct {
	RequesterID int64
	Superadmin  bool
	WorkspaceID int64
	Limit       int64
	Offset      int64
}

type CreateInput struct {
	RequesterID int64
	Superadmin  bool
	WorkspaceID int64
	Name        string
	Color       *string
}

type PatchInput struct {
	RequesterID int64
	Superadmin  bool
	WorkspaceID int64
	TagID       int64
	NameSet     bool
	Name        string
	ColorSet    bool
	ColorNull   bool
	Color       string
}

type DeleteInput struct {
	RequesterID int64
	Superadmin  bool
	WorkspaceID int64
	TagID       int64
}

func NewService(db *dbx.DB) *Service { return &Service{db: db} }

func (s *Service) List(ctx context.Context, in ListInput) ([]models.Tag, int64, error) {
	if s.db == nil {
		return nil, 0, ErrNotConfigured
	}
	if err := s.requireWorkspaceMember(ctx, in.RequesterID, in.Superadmin, in.WorkspaceID); err != nil {
		return nil, 0, err
	}
	limit, offset := normalizePaging(in.Limit, in.Offset)
	q := orm.Query[models.Tag](s.db).WhereEq("workspace_id", in.WorkspaceID)
	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	items, err := q.OrderBy("name").OrderBy("id").Limit(limit).Offset(offset).All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*models.Tag, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	if err := s.requireWorkspaceRole(ctx, in.RequesterID, in.Superadmin, in.WorkspaceID, RoleEditor); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, ErrValidation
	}
	t := &models.Tag{WorkspaceID: in.WorkspaceID, Name: name, Color: in.Color}
	if err := orm.Insert(ctx, s.db, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) Patch(ctx context.Context, in PatchInput) (*models.Tag, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	if err := s.requireWorkspaceRole(ctx, in.RequesterID, in.Superadmin, in.WorkspaceID, RoleEditor); err != nil {
		return nil, err
	}
	t, err := orm.Query[models.Tag](s.db).WhereEq("workspace_id", in.WorkspaceID).WhereEq("id", in.TagID).One(ctx)
	if err != nil {
		return nil, err
	}
	if in.NameSet {
		name := strings.TrimSpace(in.Name)
		if name == "" {
			return nil, ErrValidation
		}
		t.Name = name
	}
	if in.ColorSet {
		if in.ColorNull {
			t.Color = nil
		} else {
			v := in.Color
			t.Color = &v
		}
	}
	if err := orm.Update(ctx, s.db, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) Delete(ctx context.Context, in DeleteInput) error {
	if s.db == nil {
		return ErrNotConfigured
	}
	if err := s.requireWorkspaceRole(ctx, in.RequesterID, in.Superadmin, in.WorkspaceID, RoleEditor); err != nil {
		return err
	}
	t, err := orm.Query[models.Tag](s.db).WhereEq("workspace_id", in.WorkspaceID).WhereEq("id", in.TagID).One(ctx)
	if err != nil {
		return err
	}
	return orm.Delete(ctx, s.db, t)
}

func (s *Service) ActiveTagIDsInWorkspace(ctx context.Context, workspaceID int64, ids []int64) (map[int64]struct{}, error) {
	if len(ids) == 0 {
		return map[int64]struct{}{}, nil
	}
	rows, err := orm.Query[models.Tag](s.db).WhereEq("workspace_id", workspaceID).WhereIn("id", ids).All(ctx)
	if err != nil {
		return nil, err
	}
	out := map[int64]struct{}{}
	for _, r := range rows {
		out[r.ID] = struct{}{}
	}
	return out, nil
}

func (s *Service) FindByName(ctx context.Context, workspaceID int64, name string) (*models.Tag, error) {
	return orm.Query[models.Tag](s.db).WhereEq("workspace_id", workspaceID).UnsafeWhereExpr("LOWER(name) = LOWER({:name})", dbx.Params{"name": strings.TrimSpace(name)}).One(ctx)
}

func (s *Service) requireWorkspaceMember(ctx context.Context, requesterID int64, superadmin bool, workspaceID int64) error {
	if superadmin {
		return nil
	}
	_, err := orm.Query[models.WorkspaceUser](s.db).WhereEq("workspace_id", workspaceID).WhereEq("user_id", requesterID).One(ctx)
	if err == nil {
		return nil
	}
	if orm.HasCode(err, orm.CodeNotFound) {
		return ErrForbidden
	}
	return err
}

func (s *Service) requireWorkspaceRole(ctx context.Context, requesterID int64, superadmin bool, workspaceID int64, minRole string) error {
	if superadmin {
		return nil
	}
	m, err := orm.Query[models.WorkspaceUser](s.db).WhereEq("workspace_id", workspaceID).WhereEq("user_id", requesterID).One(ctx)
	if err != nil {
		if orm.HasCode(err, orm.CodeNotFound) {
			return ErrForbidden
		}
		return err
	}
	if roleRank(m.Role) < roleRank(minRole) {
		return ErrForbidden
	}
	return nil
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
