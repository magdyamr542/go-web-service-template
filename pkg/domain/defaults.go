package domain

import (
	"time"

	"github.com/google/uuid"
)

type DefaultFields struct {
	Id        string    `json:"id" db:"id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func NewDefaultFields() DefaultFields {
	now := time.Now()
	return DefaultFields{
		Id:        uuid.NewString(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type LimitOffset struct {
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}

const (
	MaxLimit      = 200
	DefaultLimit  = 50
	DefaultOffset = 0
)
