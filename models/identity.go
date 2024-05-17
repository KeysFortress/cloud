package models

import (
	"time"

	"github.com/google/uuid"
)

type Identity struct {
	Id         uuid.UUID
	Name       string
	KeyType    string
	KeySize    int
	PublicKey  string
	PrivateKey int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
