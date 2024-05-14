package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Password struct {
	Id                uuid.UUID
	Email             string
	Password          int
	Website           string
	CreatedAt         time.Time
	UpdatedAt         sql.NullTime
	UpperCase         bool
	LowerCase         bool
	Digits            bool
	Unique            bool
	SpecialCharacters bool
}
