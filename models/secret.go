package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Secret struct {
	Id          uuid.UUID
	Website     string
	Email       string
	Description string
	Password    int
	CreatedAt   time.Time
	UpdatedAt   sql.NullTime
}
