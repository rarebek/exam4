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


CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY,
    user_id UUID,
    post_id UUID,
    content text,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username TEXT,
    email TEXT,
    password TEXT,
    first_name TEXT,
    last_name TEXT,
    bio TEXT,
    website TEXT,
    refresh_token TEXT
);
