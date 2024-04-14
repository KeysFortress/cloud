package dtos

import "github.com/google/uuid"

type FinishAuthResponse struct {
	Uuid      uuid.UUID
	Signature string
}
