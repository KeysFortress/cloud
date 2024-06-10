package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id          uuid.UUID
	EventType   int
	Description string
	Device      sql.NullString
	EventDate   time.Time
	TypeId      int
}
