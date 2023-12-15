package domain

import (
	"time"

	"github.com/google/uuid"
)

type DefaultFields struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewDefaultFields() DefaultFields {
	now := time.Now()
	return DefaultFields{
		Id:        uuid.NewString(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}
