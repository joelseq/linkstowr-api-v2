CREATE TABLE IF NOT EXISTS links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT NOT NULL,
    title TEXT NOT NULL,
    note TEXT,
    user_id INTEGER NOT NULL,
    bookmarked_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create index on user_id in tokens table
CREATE INDEX IF NOT EXISTS idx_user_id_links ON links(user_id);

