package implementations

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	Secret string
	Issuer string
}

func (JwtService *JwtService) IssueMfaToken(id string, deviceKey string) map[string]any {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	expiresAt := time.Now().UTC().Add(time.Hour * 24 * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        id,
		"nbf":       time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp":       expiresAt,
		"iss":       JwtService.Issuer,
		"deviceKey": deviceKey,
		"role":      "mfa",
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(JwtService.Secret))

	fmt.Println(tokenString, err)
	return gin.H{"access_token": tokenString, "expires_at": expiresAt}
}

func (JwtService *JwtService) IssueToken(role string, id string, deviceKey string) map[string]any {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	expiresAt := time.Now().UTC().Add(time.Hour * 24 * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        id,
		"nbf":       time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp":       expiresAt,
		"iss":       JwtService.Issuer,
		"deviceKey": deviceKey,
		"role":      role,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(JwtService.Secret))

	fmt.Println(tokenString, err)
	return gin.H{"access_token": tokenString, "expires_at": expiresAt}
}

func (JwtService *JwtService) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(JwtService.Secret), nil
	})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["id"], claims["nbf"])

		if claims["role"] == "mfa" {
			return false, nil
		}

		return true, nil
	} else {
		fmt.Println(err)
		return false, err
	}
}

func (JwtService *JwtService) ValidateMfaToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(JwtService.Secret), nil
	})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["id"], claims["nbf"])
		return true, nil
	} else {
		fmt.Println(err)
		return false, err
	}
}

func (JwtService *JwtService) ExtractValue(tokenString string, key string) interface{} {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(JwtService.Secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["foo"], claims["nbf"])
		return claims[key]
	} else {
		fmt.Println(err)
		return ""
	}

}
