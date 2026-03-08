package models

import "time"

type User struct {
	ID           int64      `db:"id,pk"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	IsSuperadmin bool       `db:"is_superadmin"`
	Status       string     `db:"status"`
	CreatedAt    time.Time  `db:"created_at,created_at"`
	UpdatedAt    time.Time  `db:"updated_at,updated_at"`
	DeletedAt    *time.Time `db:"deleted_at,soft_delete"`
}

func (User) TableName() string { return "users" }

type UserSettings struct {
	UserID           int64     `db:"user_id,pk"`
	Theme            *string   `db:"theme,nullable"`
	Locale           *string   `db:"locale,nullable"`
	Timezone         *string   `db:"timezone,nullable"`
	BookmarksPerPage *int      `db:"bookmarks_per_page,nullable"`
	CreatedAt        time.Time `db:"created_at,created_at"`
	UpdatedAt        time.Time `db:"updated_at,updated_at"`
}

func (UserSettings) TableName() string { return "user_settings" }

type Workspace struct {
	ID          int64      `db:"id,pk"`
	Name        string     `db:"name"`
	Description *string    `db:"description,nullable"`
	OwnerUserID int64      `db:"owner_user_id"`
	CreatedAt   time.Time  `db:"created_at,created_at"`
	UpdatedAt   time.Time  `db:"updated_at,updated_at"`
	DeletedAt   *time.Time `db:"deleted_at,soft_delete"`
}

func (Workspace) TableName() string { return "workspaces" }

type WorkspaceUser struct {
	ID          int64      `db:"id,pk"`
	WorkspaceID int64      `db:"workspace_id"`
	UserID      int64      `db:"user_id"`
	Role        string     `db:"role"`
	CreatedAt   time.Time  `db:"created_at,created_at"`
	UpdatedAt   time.Time  `db:"updated_at,updated_at"`
	DeletedAt   *time.Time `db:"deleted_at,soft_delete"`
}

func (WorkspaceUser) TableName() string { return "workspace_users" }

type Bookmark struct {
	ID              int64      `db:"id,pk"`
	WorkspaceID     int64      `db:"workspace_id"`
	CreatedByUserID int64      `db:"created_by_user_id"`
	UpdatedByUserID *int64     `db:"updated_by_user_id,nullable"`
	Title           string     `db:"title"`
	URL             string     `db:"url"`
	Description     *string    `db:"description,nullable"`
	IsArchived      bool       `db:"is_archived"`
	CreatedAt       time.Time  `db:"created_at,created_at"`
	UpdatedAt       time.Time  `db:"updated_at,updated_at"`
	DeletedAt       *time.Time `db:"deleted_at,soft_delete"`
}

func (Bookmark) TableName() string { return "bookmarks" }

type Tag struct {
	ID          int64      `db:"id,pk"`
	WorkspaceID int64      `db:"workspace_id"`
	Name        string     `db:"name"`
	Color       *string    `db:"color,nullable"`
	CreatedAt   time.Time  `db:"created_at,created_at"`
	UpdatedAt   time.Time  `db:"updated_at,updated_at"`
	DeletedAt   *time.Time `db:"deleted_at,soft_delete"`
}

func (Tag) TableName() string { return "tags" }

type BookmarkTag struct {
	BookmarkID int64     `db:"bookmark_id,pk"`
	TagID      int64     `db:"tag_id,pk"`
	CreatedAt  time.Time `db:"created_at,created_at"`
}

func (BookmarkTag) TableName() string { return "bookmark_tags" }
