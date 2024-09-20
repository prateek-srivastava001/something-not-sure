package models

import (
	"time"

	"github.com/google/uuid"
)

type FriendRequest struct {
	ID            uuid.UUID `json:"id" db:"id"`
	SenderEmail   string    `json:"sender_email" db:"sender_email"`
	ReceiverEmail string    `json:"receiver_email" db:"receiver_email"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	Status        string    `json:"status" db:"status"`
}

type Friend struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserEmail   string    `json:"user_email" db:"user_email"`
	FriendEmail string    `json:"friend_email" db:"friend_email"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
