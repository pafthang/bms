CREATE TABLE IF NOT EXISTS workspace_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    workspace_id INTEGER NOT NULL REFERENCES workspaces(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    role TEXT NOT NULL CHECK (role IN ('viewer', 'editor', 'admin')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS workspace_users_workspace_user_active_uq
    ON workspace_users (workspace_id, user_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS workspace_users_user_id_idx ON workspace_users (user_id);
