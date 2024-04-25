package dtos

import "github.com/google/uuid"

type IncomingSecretsRequest struct {
	Email    string
	Password string
	Category uuid.UUID
}
