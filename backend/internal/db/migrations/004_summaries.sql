CREATE TABLE IF NOT EXISTS summaries (
    id SERIAL PRIMARY KEY,
    headline_id INT UNIQUE NOT NULL REFERENCES headlines(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    model TEXT NOT NULL DEFAULT 'gpt-4o-mini',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_summaries_headline_id ON summaries(headline_id);
