package dtos

import "github.com/google/uuid"

type IncomingPasswordUpdateRequest struct {
	Id       uuid.UUID
	Email    string
	Password string
	Website  string
}
