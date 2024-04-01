CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY,
    user_id UUID,
    post_id UUID,
    content text,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)
