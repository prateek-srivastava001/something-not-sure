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
	FirstName     string    `json:"first_name" db:"first_name"`
	LastName      string    `json:"last_name" db:"last_name"`
	Username      string    `json:"username" db:"username"`
}

type Friend struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserEmail   string    `json:"-" db:"user_email"`
	FriendEmail string    `json:"friend_email" db:"friend_email"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Username    string    `json:"username" db:"username"`
}
