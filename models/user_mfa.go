package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type UserMfa struct {
	Id     uuid.UUID
	TypeId int
	Value  sql.NullString
}
