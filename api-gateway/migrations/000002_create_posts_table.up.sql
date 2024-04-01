CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY,
    user_id UUID,
    content TEXT,
    image_url TEXT,
    title TEXT,
    likes INT,
    dislikes INT,
    views INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)
