package interfaces

import "github.com/pquerna/otp"

type TimeBasedService interface {
	GenerateTOTP(secret string, period int, algorithm otp.Algorithm) (string, error)
	GenerateOTP(secret string) (string, error)
	GenerateHOTP(secret string) (string, error)
}
