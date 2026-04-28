CREATE TABLE IF NOT EXISTS headlines (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    slug TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_headlines_created_at ON headlines(created_at DESC);
