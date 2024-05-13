package models

import (
	"time"

	"github.com/google/uuid"
)

type Password struct {
	Id                uuid.UUID
	Email             string
	Password          int
	Website           string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	UpperCase         bool
	LowerCase         bool
	Digits            bool
	Unique            bool
	SpecialCharacters bool
}
