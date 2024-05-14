package dtos

import (
	"database/sql"
	"time"
)

type CreateEvent struct {
	TypeId      int
	Description string
	DeviceId    sql.NullString
	CreatedAt   time.Time
}
