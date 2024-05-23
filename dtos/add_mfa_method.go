package dtos

import "database/sql"

type AddMfaMethod struct {
	TypeId int
	Email  sql.NullString
}
