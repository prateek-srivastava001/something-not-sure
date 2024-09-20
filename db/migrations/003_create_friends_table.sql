-- +goose Up
CREATE TABLE friends (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_email VARCHAR(100) NOT NULL,
    friend_email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_email) REFERENCES users(email) ON DELETE CASCADE,
    FOREIGN KEY (friend_email) REFERENCES users(email) ON DELETE CASCADE,
    UNIQUE (user_email, friend_email)
);

-- +goose Down
DROP TABLE IF EXISTS friends;
