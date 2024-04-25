package dtos

import "github.com/google/uuid"

type IncomingSecretsUpdateRequest struct {
	Id       uuid.UUID
	Email    string
	Password string
	Category uuid.UUID
}