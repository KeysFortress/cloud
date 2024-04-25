package dtos

import "github.com/google/uuid"

type IncomingPasswordRequest struct {
	Email    string
	Password string
	Category uuid.UUID
}
