package dtos

import "github.com/google/uuid"

type CreateIdentityResponse struct {
	Id         uuid.UUID
	PublicKey  string
	PrivateKey int
}
