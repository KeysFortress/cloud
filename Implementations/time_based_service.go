package implementations

import (
	"fmt"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
)

type TimeBasedService struct {
}

func (t *TimeBasedService) GenerateTOTP(secret string, period int, algorithm otp.Algorithm) (string, error) {

	code, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Algorithm: algorithm,
		Period:    uint(period),
	})
	if err != nil {
		fmt.Println("Error generating code:", err)
		return "", err
	}

	fmt.Println("Current TOTP code:", code)
	return code, nil
}

func (t *TimeBasedService) GenerateHOTP(secret string) ([]string, error) {

	key, err := hotp.Generate(hotp.GenerateOpts{
		Issuer:      "Example Corp",
		AccountName: "user@example.com",
		Secret:      []byte(secret),
	})
	if err != nil {
		fmt.Println("Error generating HOTP:", err)
		return []string{}, err
	}

	var codes []string

	for i := 0; i < 5; i++ {
		code, err := t.GenerateTOTP(key.Secret(), int(key.Period()), key.Algorithm())
		if err != nil {
			fmt.Println("Failed to generate code for HOTP")
			return []string{}, nil
		}
		fmt.Println(code)
		codes = append(codes, code)
	}
	return codes, nil
}
