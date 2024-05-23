package interfaces

import "github.com/pquerna/otp"

type TimeBasedService interface {
	GenerateTOTPCode(secret string, period int, algorithm otp.Algorithm) (string, error)
	GenerateHOTPCode(secret string) ([]string, error)
	VerifyTOTP(code string, secret string, period int, algorithm otp.Algorithm) (bool, error)
	GenerateTOTP(accountName string) (string, error)
}
