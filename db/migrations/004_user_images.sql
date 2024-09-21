-- +goose Up
CREATE TABLE user_media (
    id UUID PRIMARY KEY NOT NULL,
    user_email VARCHAR(100) NOT NULL,
    image_url TEXT,
    audio_url TEXT,
    image_parsed TEXT,
    audio_parsed TEXT,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_email) REFERENCES users(email) ON DELETE CASCADE,
    UNIQUE (user_email, image_url, audio_url)
);

-- +goose Down
DROP TABLE IF EXISTS user_media;
