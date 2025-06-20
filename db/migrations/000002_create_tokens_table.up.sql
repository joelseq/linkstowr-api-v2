CREATE TABLE IF NOT EXISTS tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    token_hash TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    short_token TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create index on token_hash in tokens table
CREATE INDEX IF NOT EXISTS idx_token_hash_tokens ON tokens(token_hash);

-- Create index on user_id in tokens table
CREATE INDEX IF NOT EXISTS idx_user_id_tokens ON tokens(user_id);

