package models

import (
	"database/sql"
	"time"
)

type TimeBasedCode struct {
	Id        string
	Email     string
	Website   string
	Code      string
	Secret    int
	Type      string
	Validity  int
	Algorithm string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
