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

type IdentityInternal struct {
	Id         uuid.UUID
	Name       string
	KeyType    int
	KeySize    int
	PublicKey  string
	PrivateKey string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
