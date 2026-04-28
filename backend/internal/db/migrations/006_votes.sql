CREATE TABLE IF NOT EXISTS votes (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    headline_id INT NOT NULL REFERENCES headlines(id) ON DELETE CASCADE,
    value INT NOT NULL CHECK (value IN (-1, 1)),
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, headline_id)
);

CREATE INDEX IF NOT EXISTS idx_votes_headline_id ON votes(headline_id);
