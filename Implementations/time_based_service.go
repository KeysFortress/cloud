package implementations

import (
	"fmt"
	"log"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
)

type TimeBasedService struct {
	Issuer string
}

func (t *TimeBasedService) GenerateTOTPCode(secret string, period int, algorithm otp.Algorithm) (string, error) {

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

func (t *TimeBasedService) GenerateHOTPCode(secret string) ([]string, error) {

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
		code, err := t.GenerateTOTPCode(key.Secret(), int(key.Period()), key.Algorithm())
		if err != nil {
			fmt.Println("Failed to generate code for HOTP")
			return []string{}, nil
		}
		fmt.Println(code)
		codes = append(codes, code)
	}
	return codes, nil
}

func (t *TimeBasedService) VerifyTOTP(code string, secret string, period int, algorithm otp.Algorithm) (bool, error) {
	isValid, err := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period:    uint(period),
		Algorithm: algorithm,
	})

	if err != nil {
		return false, err
	}

	return isValid, nil
}

func (t *TimeBasedService) GenerateTOTP(accountName string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      t.Issuer,
		AccountName: accountName,
	})

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	fmt.Println("Key URL:", key.URL())

	fmt.Println("Secret:", key.Secret())

	return key.Secret(), nil
}
