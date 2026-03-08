package settings

import (
	"context"
	"errors"

	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/domain/models"
	"github.com/pafthang/dbx"
	"github.com/pafthang/orm"
)

var ErrNotConfigured = errors.New("database is not configured")

type Service struct {
	db *dbx.DB
}

type PutInput struct {
	Theme            *string
	Locale           *string
	Timezone         *string
	BookmarksPerPage *int
}

type PatchInput struct {
	Theme            arc.OptionalString
	Locale           arc.OptionalString
	Timezone         arc.OptionalString
	BookmarksPerPage arc.OptionalInt
}

func NewService(db *dbx.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Get(ctx context.Context, userID int64) (*models.UserSettings, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	settings, err := orm.ByPK[models.UserSettings](ctx, s.db, userID)
	if err == nil {
		return settings, nil
	}
	if !orm.HasCode(err, orm.CodeNotFound) {
		return nil, err
	}
	fresh := &models.UserSettings{UserID: userID}
	if err := orm.Insert(ctx, s.db, fresh); err != nil {
		if orm.HasCode(err, orm.CodeConflict) {
			return orm.ByPK[models.UserSettings](ctx, s.db, userID)
		}
		return nil, err
	}
	return fresh, nil
}

func (s *Service) Put(ctx context.Context, userID int64, in PutInput) (*models.UserSettings, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	settings, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	settings.Theme = in.Theme
	settings.Locale = in.Locale
	settings.Timezone = in.Timezone
	settings.BookmarksPerPage = in.BookmarksPerPage
	if err := orm.Update(ctx, s.db, settings); err != nil {
		return nil, err
	}
	return settings, nil
}

func (s *Service) Patch(ctx context.Context, userID int64, in PatchInput) (*models.UserSettings, error) {
	if s.db == nil {
		return nil, ErrNotConfigured
	}
	settings, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if in.Theme.IsSet() {
		if in.Theme.IsNull() {
			settings.Theme = nil
		} else if v, ok := in.Theme.Value(); ok {
			settings.Theme = &v
		}
	}
	if in.Locale.IsSet() {
		if in.Locale.IsNull() {
			settings.Locale = nil
		} else if v, ok := in.Locale.Value(); ok {
			settings.Locale = &v
		}
	}
	if in.Timezone.IsSet() {
		if in.Timezone.IsNull() {
			settings.Timezone = nil
		} else if v, ok := in.Timezone.Value(); ok {
			settings.Timezone = &v
		}
	}
	if in.BookmarksPerPage.IsSet() {
		if in.BookmarksPerPage.IsNull() {
			settings.BookmarksPerPage = nil
		} else if v, ok := in.BookmarksPerPage.Value(); ok {
			settings.BookmarksPerPage = &v
		}
	}
	if err := orm.Update(ctx, s.db, settings); err != nil {
		return nil, err
	}
	return settings, nil
}
