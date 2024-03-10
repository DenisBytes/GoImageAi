package types

import (
	"time"

	"github.com/google/uuid"
)

type ImageStatus int

const (
	ImageStatusFailed ImageStatus = iota
	ImageStatusPending
	ImageStatusCompleted
)

type Image struct {
	ID            int `bun:"id,pk,autoincrement"`
	UserID        uuid.UUID
	Status        ImageStatus
	Prompt        string
	ImageLocation string
	Deleted       bool
	CreatedAt     time.Time `bun:"default:'now()'"`
	DeletedAt     time.Time
}
