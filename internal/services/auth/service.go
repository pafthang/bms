package auth

import (
	"context"
	"errors"
	"strings"

	bmsauth "github.com/pafthang/bms/internal/auth"
	"github.com/pafthang/bms/internal/domain/models"
	"github.com/pafthang/dbx"
	"github.com/pafthang/orm"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrNotConfigured = errors.New("database is not configured")

type Service struct {
	db    *dbx.DB
	token *bmsauth.JWTManager
}

type UserView struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	IsSuperadmin bool   `json:"is_superadmin"`
	Status       string `json:"status"`
}

func NewService(db *dbx.DB, token *bmsauth.JWTManager) *Service {
	return &Service{db: db, token: token}
}

func (s *Service) Register(ctx context.Context, email, password string) (*UserView, *bmsauth.TokenPair, error) {
	if s.db == nil {
		return nil, nil, ErrNotConfigured
	}
	normEmail := normalizeEmail(email)
	hash, err := bmsauth.HashPassword(password)
	if err != nil {
		return nil, nil, err
	}
	user := &models.User{
		Email:        normEmail,
		PasswordHash: hash,
		Status:       "active",
	}
	if err := orm.WithTx(ctx, s.db, func(tx *orm.Tx) error {
		if err := orm.Insert(ctx, tx, user); err != nil {
			return err
		}
		settings := &models.UserSettings{UserID: user.ID}
		if err := orm.Insert(ctx, tx, settings); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, nil, err
	}
	pair, err := s.token.IssuePair(user.ID, user.Email, user.IsSuperadmin)
	if err != nil {
		return nil, nil, err
	}
	return toUserView(user), pair, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*UserView, *bmsauth.TokenPair, error) {
	if s.db == nil {
		return nil, nil, ErrNotConfigured
	}
	normEmail := normalizeEmail(email)
	user, err := orm.Query[models.User](s.db).WhereEq("email", normEmail).One(ctx)
	if err != nil {
		if orm.HasCode(err, orm.CodeNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}
	if err := bmsauth.CheckPassword(user.PasswordHash, password); err != nil {
		return nil, nil, ErrInvalidCredentials
	}
	pair, err := s.token.IssuePair(user.ID, user.Email, user.IsSuperadmin)
	if err != nil {
		return nil, nil, err
	}
	return toUserView(user), pair, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*bmsauth.TokenPair, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	claims, err := s.token.VerifyRefresh(strings.TrimSpace(refreshToken))
	if err != nil {
		return nil, err
	}
	user, err := orm.ByPK[models.User](ctx, s.db, claims.Subject)
	if err != nil {
		return nil, err
	}
	return s.token.IssuePair(user.ID, user.Email, user.IsSuperadmin)
}

func (s *Service) Me(ctx context.Context, userID int64) (*UserView, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	user, err := orm.ByPK[models.User](ctx, s.db, userID)
	if err != nil {
		return nil, err
	}
	return toUserView(user), nil
}

func (s *Service) VerifyAccess(token string) (*bmsauth.Claims, error) {
	return s.token.VerifyAccess(token)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func toUserView(u *models.User) *UserView {
	return &UserView{ID: u.ID, Email: u.Email, IsSuperadmin: u.IsSuperadmin, Status: u.Status}
}
