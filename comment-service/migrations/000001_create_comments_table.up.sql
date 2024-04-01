CREATE TABLE IF NOT EXISTS comments (
    id uuid,
    content TEXT,
    post_id uuid,
    owner_id uuid,
    created_at TIME DEFAULT NOW(),
    updated_at TIME,
    created_at TIME
);