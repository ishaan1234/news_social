CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    headline_id INT REFERENCES headlines(id) ON DELETE CASCADE,
    source TEXT,
    url TEXT UNIQUE,
    content TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);