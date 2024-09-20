package models

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	CreatorEmail string    `json:"creator_email" db:"creator_email"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type GroupMember struct {
	ID          uuid.UUID `json:"id" db:"id"`
	GroupID     uuid.UUID `json:"group_id" db:"group_id"`
	MemberEmail string    `json:"member_email" db:"member_email"`
	JoinedAt    time.Time `json:"joined_at" db:"joined_at"`
}
