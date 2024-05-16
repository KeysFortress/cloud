package dtos

import "github.com/google/uuid"

type UpdateTimeBasedCode struct {
	Id        uuid.UUID
	Website   string
	Email     string
	Secret    string
	Type      int
	Validity  int
	Algorithm int
}
