INSERT OR IGNORE INTO users (email, password_hash, is_superadmin, status)
VALUES ('admin@admin.admin', '$2a$10$YgLNVxzfzZcrcKtGN5FmWOMa/N/v3a/l7QNc5wffJ.z20FNv/ciq.', TRUE, 'active');

INSERT OR IGNORE INTO user_settings (user_id)
SELECT id
FROM users
WHERE email = 'admin@admin.admin';
