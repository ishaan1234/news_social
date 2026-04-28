CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    headline_id INT NOT NULL REFERENCES headlines(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_comments_headline_id ON comments(headline_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
