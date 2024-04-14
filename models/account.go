package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
	Enabled   bool
}
