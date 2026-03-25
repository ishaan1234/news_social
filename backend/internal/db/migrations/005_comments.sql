CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    headline_id INT REFERENCES headlines(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);