-- +goose Up
CREATE TABLE friend_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_email VARCHAR(100) NOT NULL,
    receiver_email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status TEXT CHECK (status IN ('pending', 'accepted', 'rejected')) DEFAULT 'pending',
    FOREIGN KEY (sender_email) REFERENCES users(email) ON DELETE CASCADE,
    FOREIGN KEY (receiver_email) REFERENCES users(email) ON DELETE CASCADE,
    UNIQUE (sender_email, receiver_email) -- Prevent duplicate requests
);

-- +goose Down
DROP TABLE IF EXISTS friend_requests;
