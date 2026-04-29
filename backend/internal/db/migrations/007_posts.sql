CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    author_id TEXT,
    author_name TEXT NOT NULL DEFAULT 'Anonymous',
    author_handle TEXT,
    body TEXT NOT NULL,
    article_url TEXT NOT NULL,
    article_title TEXT NOT NULL,
    article_source TEXT,
    article_summary TEXT,
    article_image_url TEXT,
    article_published_at TEXT,
    share_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_article_url ON posts(article_url);

CREATE TABLE IF NOT EXISTS post_comments (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    author_id TEXT,
    author_name TEXT NOT NULL DEFAULT 'Anonymous',
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_post_comments_post_id ON post_comments(post_id);

CREATE TABLE IF NOT EXISTS post_votes (
    post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    voter_id TEXT NOT NULL,
    value INT NOT NULL CHECK (value IN (-1, 1)),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (post_id, voter_id)
);

CREATE INDEX IF NOT EXISTS idx_post_votes_post_id ON post_votes(post_id);
