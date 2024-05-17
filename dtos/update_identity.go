package dtos

import "github.com/google/uuid"

type UpdateIdentity struct {
	Id            uuid.UUID
	Name          string
	KeyType       int
	KeySize       int
	RegenerateKey bool
}
