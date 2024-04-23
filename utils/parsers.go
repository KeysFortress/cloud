package utils

import (
	"github.com/google/uuid"
)

func ParseUUID(in string) (uuid.UUID, error) {
	id, err := uuid.Parse(in)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
