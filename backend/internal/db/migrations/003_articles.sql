CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    headline_id INT NOT NULL REFERENCES headlines(id) ON DELETE CASCADE,
    source TEXT,
    title TEXT,
    url TEXT UNIQUE,
    content TEXT,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_articles_headline_id ON articles(headline_id);
