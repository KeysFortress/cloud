package dtos

import (
	"time"

	"github.com/google/uuid"
)

type StoredAuthRequest struct {
	Name string
	Code string
	Uuid string
	Id   uuid.UUID
	Time time.Time
}
