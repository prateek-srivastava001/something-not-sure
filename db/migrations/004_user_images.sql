-- +goose Up
CREATE TABLE user_media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_email VARCHAR(100) NOT NULL,
    image_url TEXT,
    audio_url TEXT,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_email) REFERENCES users(email) ON DELETE CASCADE,
    UNIQUE (user_email, image_url, audio_url)
);

-- +goose Down
DROP TABLE IF EXISTS user_media;
