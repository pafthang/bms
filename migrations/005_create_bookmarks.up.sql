CREATE TABLE IF NOT EXISTS bookmarks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    workspace_id INTEGER NOT NULL REFERENCES workspaces(id),
    created_by_user_id INTEGER NOT NULL REFERENCES users(id),
    updated_by_user_id INTEGER NULL REFERENCES users(id),
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    description TEXT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS bookmarks_workspace_created_idx ON bookmarks (workspace_id, created_at DESC);
CREATE INDEX IF NOT EXISTS bookmarks_workspace_archived_idx ON bookmarks (workspace_id, is_archived);
CREATE INDEX IF NOT EXISTS bookmarks_workspace_deleted_idx ON bookmarks (workspace_id, deleted_at);
