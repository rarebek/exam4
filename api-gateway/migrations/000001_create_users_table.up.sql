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
