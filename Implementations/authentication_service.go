package implementations

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/google/uuid"

	"leanmeal/api/dtos"
)

type AuthenticationService struct {
	AuthRequests map[uuid.UUID]dtos.StoredAuthRequest
}

func (authService *AuthenticationService) GetMessage(email *string, id *uuid.UUID) dtos.InitAuthReponse {

	if authService.AuthRequests == nil {
		authService.AuthRequests = make(map[uuid.UUID]dtos.StoredAuthRequest)
	}

	uuid := uuid.New()
	code := generateRandomString(32)

	authResponse := dtos.StoredAuthRequest{
		Id:       *id,
		Name:     *email,
		Code:     code,
		Uuid:     uuid.String(),
		Time:     time.Time{}.Add(time.Duration(time.Second * 30)),
		Approved: false,
	}

	authService.AuthRequests[uuid] = authResponse

	fmt.Println(&authService.AuthRequests)

	response := dtos.InitAuthReponse{
		Code: authResponse.Code,
		Uuid: authResponse.Uuid,
	}

	return response
}

func (authService *AuthenticationService) GetRequestById(id uuid.UUID) uuid.UUID {
	return authService.AuthRequests[id].Id
}

func (authService *AuthenticationService) VerifySignature(response dtos.FinishAuthResponse, keys *[]string) (uuid.UUID, error) {
	authRequest := authService.AuthRequests[response.Uuid]
	messege := authRequest.Code
	for _, key := range *keys {
		pk, err := base64.StdEncoding.DecodeString(key)

		if err != nil {
			fmt.Println("Failed to decode public key")
			fmt.Println(err)
			return uuid.UUID{}, err
		}

		decodedSignature, err := base64.StdEncoding.DecodeString(response.Signature)

		if err != nil {
			fmt.Println("Failed to decode Signature")
			fmt.Println(err)
			return uuid.UUID{}, err
		}

		isValid := ed25519.Verify(pk, []byte(messege), decodedSignature)
		if isValid {
			authRequest.Approved = true
			return authRequest.Id, nil
		}
	}

	return uuid.Nil, nil
}

func (authService *AuthenticationService) ExchangeCodeForToken(code uuid.UUID) (uuid.UUID, bool) {
	authRequest, ok := authService.AuthRequests[code]

	if !ok {
		return uuid.UUID{}, false
	}

	return authRequest.Id, true
}

func (authService *AuthenticationService) Start() {
	// Create a channel to receive signals
	signal := make(chan struct{})

	// Start a goroutine to send signals at regular intervals
	go func() {
		for {
			time.Sleep(30 * time.Second) // Wait for 30 seconds
			signal <- struct{}{}         // Send a signal to the channel
		}
	}()

	for {
		select {
		case <-signal:
			fmt.Println("Something has happened at", time.Now())
			for i, d := range authService.AuthRequests {
				fmt.Println(i)
				fmt.Println("Checking if request has expired")
				expiresAt := d.Time.UTC()
				fmt.Fprintln(os.Stdout, "Expires at ", expiresAt, "and the time now is", time.Now().UTC())
				expired := d.Time.UTC().Before(time.Now().UTC())
				fmt.Println(expired)
				if expired {
					delete(authService.AuthRequests, i)
					fmt.Println("Timer has expired")
					fmt.Println(d.Uuid)
					fmt.Println(d.Time)
				}
			}
			fmt.Println(authService.AuthRequests)
		}
	}
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string
	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result += string(charset[randomIndex.Int64()])
	}
	return result
}
