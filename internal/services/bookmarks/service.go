package bookmarks

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/pafthang/bms/internal/domain/models"
	"github.com/pafthang/bms/internal/services/tags"
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
	db   *dbx.DB
	tags *tags.Service
}

type ListInput struct {
	RequesterID int64
	Superadmin  bool
	WorkspaceID int64
	Limit       int64
	Offset      int64
	Search      string
	Tag         string
	TagID       *int64
	IsArchived  *bool
	Sort        string
	Order       string
}

type CreateInput struct {
	RequesterID int64
	Superadmin  bool
	WorkspaceID int64
	Title       string
	URL         string
	Description *string
	IsArchived  *bool
	TagIDs      []int64
}

type PatchInput struct {
	RequesterID     int64
	Superadmin      bool
	WorkspaceID     int64
	BookmarkID      int64
	TitleSet        bool
	Title           string
	URLSet          bool
	URL             string
	DescriptionSet  bool
	DescriptionNull bool
	Description     string
	IsArchivedSet   bool
	IsArchived      bool
	TagIDsSet       bool
	TagIDs          []int64
	UpdatedByUserID int64
}

type DeleteInput struct {
	RequesterID int64
	Superadmin  bool
	WorkspaceID int64
	BookmarkID  int64
}

type BookmarkView struct {
	Model  models.Bookmark
	TagIDs []int64
}

func NewService(db *dbx.DB, tagsSvc *tags.Service) *Service {
	return &Service{db: db, tags: tagsSvc}
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*BookmarkView, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	if err := s.requireWorkspaceRole(ctx, in.RequesterID, in.Superadmin, in.WorkspaceID, RoleEditor); err != nil {
		return nil, err
	}
	title := strings.TrimSpace(in.Title)
	if title == "" || !isValidHTTPURL(in.URL) {
		return nil, ErrValidation
	}
	archived := false
	if in.IsArchived != nil {
		archived = *in.IsArchived
	}
	bookmark := &models.Bookmark{
		WorkspaceID:     in.WorkspaceID,
		CreatedByUserID: in.RequesterID,
		Title:           title,
		URL:             strings.TrimSpace(in.URL),
		Description:     in.Description,
		IsArchived:      archived,
	}
	if err := orm.WithTx(ctx, s.db, func(tx *orm.Tx) error {
		if err := orm.Insert(ctx, tx, bookmark); err != nil {
			return err
		}
		return s.replaceBookmarkTagsTx(ctx, tx, in.WorkspaceID, bookmark.ID, in.TagIDs)
	}); err != nil {
		return nil, err
	}
	return s.Get(ctx, in.RequesterID, in.Superadmin, in.WorkspaceID, bookmark.ID)
}

func (s *Service) List(ctx context.Context, in ListInput) ([]BookmarkView, int64, error) {
	if s.db == nil {
		return nil, 0, ErrNotConfigured
	}
	if err := s.requireWorkspaceMember(ctx, in.RequesterID, in.Superadmin, in.WorkspaceID); err != nil {
		return nil, 0, err
	}
	limit, offset := normalizePaging(in.Limit, in.Offset)
	q := orm.Query[models.Bookmark](s.db).WhereEq("workspace_id", in.WorkspaceID)
	if in.IsArchived != nil {
		q = q.WhereEq("is_archived", *in.IsArchived)
	}
	if strings.TrimSpace(in.Search) != "" {
		term := "%" + strings.TrimSpace(in.Search) + "%"
		q = q.UnsafeWhereExpr("(title LIKE {:q} OR url LIKE {:q} OR description LIKE {:q})", dbx.Params{"q": term})
	}
	if err := s.applyTagFilters(ctx, q, in.WorkspaceID, strings.TrimSpace(in.Tag), in.TagID); err != nil {
		if errors.Is(err, orm.ErrNotFound) {
			return []BookmarkView{}, 0, nil
		}
		return nil, 0, err
	}
	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	q = applySort(q, in.Sort, in.Order)
	items, err := q.Limit(limit).Offset(offset).All(ctx)
	if err != nil {
		return nil, 0, err
	}
	out := make([]BookmarkView, 0, len(items))
	for _, b := range items {
		tagIDs, err := s.fetchTagIDs(ctx, b.ID)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, BookmarkView{Model: b, TagIDs: tagIDs})
	}
	return out, total, nil
}

func (s *Service) Get(ctx context.Context, requesterID int64, superadmin bool, workspaceID, bookmarkID int64) (*BookmarkView, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	if err := s.requireWorkspaceMember(ctx, requesterID, superadmin, workspaceID); err != nil {
		return nil, err
	}
	b, err := orm.Query[models.Bookmark](s.db).WhereEq("workspace_id", workspaceID).WhereEq("id", bookmarkID).One(ctx)
	if err != nil {
		return nil, err
	}
	tagIDs, err := s.fetchTagIDs(ctx, b.ID)
	if err != nil {
		return nil, err
	}
	return &BookmarkView{Model: *b, TagIDs: tagIDs}, nil
}

func (s *Service) Patch(ctx context.Context, in PatchInput) (*BookmarkView, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	if err := s.requireWorkspaceRole(ctx, in.RequesterID, in.Superadmin, in.WorkspaceID, RoleEditor); err != nil {
		return nil, err
	}
	b, err := orm.Query[models.Bookmark](s.db).WhereEq("workspace_id", in.WorkspaceID).WhereEq("id", in.BookmarkID).One(ctx)
	if err != nil {
		return nil, err
	}
	if in.TitleSet {
		title := strings.TrimSpace(in.Title)
		if title == "" {
			return nil, ErrValidation
		}
		b.Title = title
	}
	if in.URLSet {
		if !isValidHTTPURL(in.URL) {
			return nil, ErrValidation
		}
		b.URL = strings.TrimSpace(in.URL)
	}
	if in.DescriptionSet {
		if in.DescriptionNull {
			b.Description = nil
		} else {
			v := in.Description
			b.Description = &v
		}
	}
	if in.IsArchivedSet {
		b.IsArchived = in.IsArchived
	}
	uid := in.UpdatedByUserID
	b.UpdatedByUserID = &uid
	if err := orm.WithTx(ctx, s.db, func(tx *orm.Tx) error {
		if err := orm.Update(ctx, tx, b); err != nil {
			return err
		}
		if in.TagIDsSet {
			if err := s.replaceBookmarkTagsTx(ctx, tx, in.WorkspaceID, b.ID, in.TagIDs); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return s.Get(ctx, in.RequesterID, in.Superadmin, in.WorkspaceID, b.ID)
}

func (s *Service) Delete(ctx context.Context, in DeleteInput) error {
	if s.db == nil {
		return ErrNotConfigured
	}
	if err := s.requireWorkspaceRole(ctx, in.RequesterID, in.Superadmin, in.WorkspaceID, RoleEditor); err != nil {
		return err
	}
	b, err := orm.Query[models.Bookmark](s.db).WhereEq("workspace_id", in.WorkspaceID).WhereEq("id", in.BookmarkID).One(ctx)
	if err != nil {
		return err
	}
	return orm.Delete(ctx, s.db, b)
}

func (s *Service) replaceBookmarkTagsTx(ctx context.Context, tx *orm.Tx, workspaceID, bookmarkID int64, tagIDs []int64) error {
	if _, err := orm.Query[models.BookmarkTag](tx).WhereEq("bookmark_id", bookmarkID).Delete(ctx); err != nil {
		return err
	}
	if len(tagIDs) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(tagIDs))
	uniq := make([]int64, 0, len(tagIDs))
	for _, id := range tagIDs {
		if id <= 0 {
			return ErrValidation
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniq = append(uniq, id)
	}
	activeIDs, err := s.tags.ActiveTagIDsInWorkspace(ctx, workspaceID, uniq)
	if err != nil {
		return err
	}
	if len(activeIDs) != len(uniq) {
		return ErrValidation
	}
	for _, tagID := range uniq {
		bt := &models.BookmarkTag{BookmarkID: bookmarkID, TagID: tagID}
		if err := orm.Insert(ctx, tx, bt); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) fetchTagIDs(ctx context.Context, bookmarkID int64) ([]int64, error) {
	rows, err := orm.Query[models.BookmarkTag](s.db).WhereEq("bookmark_id", bookmarkID).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]int64, 0, len(rows))
	for _, r := range rows {
		_, err := orm.ByPK[models.Tag](ctx, s.db, r.TagID)
		if err == nil {
			out = append(out, r.TagID)
			continue
		}
		if orm.HasCode(err, orm.CodeNotFound) || orm.HasCode(err, orm.CodeSoftDeleted) {
			continue
		}
		return nil, err
	}
	return out, nil
}

func (s *Service) applyTagFilters(ctx context.Context, q *orm.ModelQuery[models.Bookmark], workspaceID int64, tagName string, tagID *int64) error {
	var filterTagIDs []int64
	if tagID != nil && *tagID > 0 {
		filterTagIDs = append(filterTagIDs, *tagID)
	}
	if tagName != "" {
		t, err := s.tags.FindByName(ctx, workspaceID, tagName)
		if err != nil {
			return err
		}
		filterTagIDs = append(filterTagIDs, t.ID)
	}
	if len(filterTagIDs) == 0 {
		return nil
	}
	links, err := orm.Query[models.BookmarkTag](s.db).WhereIn("tag_id", filterTagIDs).All(ctx)
	if err != nil {
		return err
	}
	if len(links) == 0 {
		return orm.ErrNotFound
	}
	bookmarkIDs := make([]int64, 0, len(links))
	uniq := map[int64]struct{}{}
	for _, l := range links {
		if _, ok := uniq[l.BookmarkID]; ok {
			continue
		}
		uniq[l.BookmarkID] = struct{}{}
		bookmarkIDs = append(bookmarkIDs, l.BookmarkID)
	}
	q.WhereIn("id", bookmarkIDs)
	return nil
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

func applySort(q *orm.ModelQuery[models.Bookmark], sortField, order string) *orm.ModelQuery[models.Bookmark] {
	field := strings.TrimSpace(sortField)
	if field == "" {
		field = "created_at"
	}
	switch field {
	case "created_at", "updated_at", "title":
	default:
		field = "created_at"
	}
	desc := strings.EqualFold(strings.TrimSpace(order), "desc") || strings.TrimSpace(order) == ""
	if desc {
		return q.OrderByDesc(field).OrderByDesc("id")
	}
	return q.OrderBy(field).OrderBy("id")
}

func isValidHTTPURL(v string) bool {
	u, err := url.ParseRequestURI(strings.TrimSpace(v))
	if err != nil {
		return false
	}
	scheme := strings.ToLower(strings.TrimSpace(u.Scheme))
	return (scheme == "http" || scheme == "https") && u.Host != ""
}
