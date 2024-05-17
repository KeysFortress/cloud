package utils

import (
	"github.com/google/uuid"
	"github.com/pquerna/otp"
)

func ParseUUID(in string) (uuid.UUID, error) {
	id, err := uuid.Parse(in)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func ParseAlgorithm(current int) otp.Algorithm {
	var algorithm otp.Algorithm

	switch current {
	case 1:
		algorithm = otp.AlgorithmSHA1
	case 2:
		algorithm = otp.AlgorithmSHA256
	case 3:
		algorithm = otp.AlgorithmSHA512
	default:
		algorithm = otp.AlgorithmMD5
	}

	return algorithm
}
