package users

import (
	"context"
	"errors"
	"strings"

	"github.com/pafthang/bms/internal/domain/models"
	"github.com/pafthang/dbx"
	"github.com/pafthang/orm"
)

var ErrNotConfigured = errors.New("database is not configured")
var ErrEmailChangeForbidden = errors.New("email change is forbidden in mvp")

type PatchMeInput struct {
	EmailSet bool
	Email    string
}

type Service struct {
	db *dbx.DB
}

func NewService(db *dbx.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Me(ctx context.Context, userID int64) (*models.User, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	return orm.ByPK[models.User](ctx, s.db, userID)
}

func (s *Service) PatchMe(ctx context.Context, userID int64, in PatchMeInput) (*models.User, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	user, err := orm.ByPK[models.User](ctx, s.db, userID)
	if err != nil {
		return nil, err
	}
	if in.EmailSet {
		normalized := strings.ToLower(strings.TrimSpace(in.Email))
		if normalized != user.Email {
			return nil, ErrEmailChangeForbidden
		}
	}
	return user, nil
}

func (s *Service) AdminList(ctx context.Context, limit, offset int64, search string) ([]models.User, int64, error) {
	if s.db == nil {
		return nil, 0, ErrNotConfigured
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	q := orm.Query[models.User](s.db)
	if strings.TrimSpace(search) != "" {
		term := "%" + strings.TrimSpace(search) + "%"
		q = q.UnsafeWhereExpr("(email LIKE {:q} OR status LIKE {:q})", dbx.Params{"q": term})
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

func (s *Service) AdminGet(ctx context.Context, userID int64) (*models.User, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	return orm.ByPK[models.User](ctx, s.db, userID)
}
