CREATE TABLE IF NOT EXISTS summaries (
    id SERIAL PRIMARY KEY,
    headline_id INT UNIQUE REFERENCES headlines(id) ON DELETE CASCADE,
    summary TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);