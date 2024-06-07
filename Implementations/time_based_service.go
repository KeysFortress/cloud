package implementations

import (
	"fmt"
	"log"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"

	"leanmeal/api/dtos"
)

type TimeBasedService struct {
	Issuer string
}

func (t *TimeBasedService) GenerateTOTPCode(secret string, period int, algorithm otp.Algorithm) (string, error) {
	code, err := totp.GenerateCodeCustom(secret, time.Now().UTC(), totp.ValidateOpts{
		Algorithm: algorithm,
		Period:    uint(period),
		Digits:    6,
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
		Issuer:      t.Issuer,
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
	isValid, _ := totp.ValidateCustom(code, secret, time.Now().UTC(), totp.ValidateOpts{
		Period:    uint(period),
		Digits:    6,
		Algorithm: algorithm,
	})

	return isValid, nil
}

func (t *TimeBasedService) GenerateTOTP(accountName string, period int, algorithm otp.Algorithm) (dtos.MfaSetupResponse, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      t.Issuer,
		AccountName: accountName,
		Algorithm:   algorithm,
		Digits:      6,
		Period:      uint(30),
	})

	if err != nil {
		log.Fatal(err)
		return dtos.MfaSetupResponse{}, err
	}

	fmt.Println("Key URL:", key.URL())

	fmt.Println("Secret:", key.Secret())

	code, _ := t.GenerateTOTPCode(key.Secret(), 30, algorithm)
	isValid, _ := t.VerifyTOTP(code, key.Secret(), 30, algorithm)

	fmt.Println(time.Now().UTC())

	fmt.Println("Code is", code, " it's ", isValid, " secret = ", key.Secret())

	return dtos.MfaSetupResponse{
		Secret: key.Secret(),
		QrCode: key.URL(),
	}, nil
}
