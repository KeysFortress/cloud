package interfaces

import (
	"github.com/pquerna/otp"

	"leanmeal/api/dtos"
)

type TimeBasedService interface {
	GenerateTOTPCode(secret string, period int, algorithm otp.Algorithm) (string, error)
	GenerateHOTPCode(secret string) ([]string, error)
	VerifyTOTP(code string, secret string, period int, algorithm otp.Algorithm) (bool, error)
	GenerateTOTP(accountName string, period int, algorithm otp.Algorithm) (dtos.MfaSetupResponse, error)
}
