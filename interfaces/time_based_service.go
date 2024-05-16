package interfaces

import "github.com/pquerna/otp"

type TimeBasedService interface {
	GenerateTOTP(secret string, period int, algorithm otp.Algorithm) (string, error)
	GenerateHOTP(secret string) (string, error)
	VerifyTOTP(code string, secret string, period int, algorithm otp.Algorithm) (bool, error)
}
