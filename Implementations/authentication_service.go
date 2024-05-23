package implementations

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/syncmap"

	"leanmeal/api/dtos"
)

type AuthenticationService struct {
	AuthRequests syncmap.Map
	Domain       string
}

func (authService *AuthenticationService) GetMessage(email *string, id *uuid.UUID) dtos.InitAuthReponse {

	uuid := uuid.New()
	code := generateRandomString(32)

	authResponse := dtos.StoredAuthRequest{
		Id:       *id,
		Name:     *email,
		Code:     code,
		Uuid:     uuid.String(),
		Time:     time.Now().UTC().Add(time.Duration(time.Minute * 10)),
		Approved: false,
	}

	authService.AuthRequests.LoadOrStore(uuid, authResponse)

	fmt.Println(&authService.AuthRequests)

	response := dtos.InitAuthReponse{
		Code:            authResponse.Code,
		Uuid:            authResponse.Uuid,
		Domain:          authService.Domain,
		VerifySignature: authService.Domain + "/v1/finish-request",
	}

	return response
}

func (authService *AuthenticationService) GetRequestById(id uuid.UUID) (uuid.UUID, error) {
	request, exists := authService.AuthRequests.Load(id)
	if !exists {
		return uuid.UUID{}, errors.New("record doesn't exist")
	}
	return request.(dtos.StoredAuthRequest).Id, nil
}

func (authService *AuthenticationService) VerifySignature(response dtos.FinishAuthResponse, keys *[]string) (uuid.UUID, error) {
	request, _ := authService.AuthRequests.Load(response.Uuid)

	authRequest := request.(dtos.StoredAuthRequest)
	currentCode := authRequest.Code

	if currentCode == "" {
		return uuid.Nil, errors.New("failed to get the challange, expired")
	}

	for _, key := range *keys {
		pk, err := base64.StdEncoding.DecodeString(key)

		if err != nil {
			fmt.Println("Failed to decode public key")
			fmt.Println(err)
			return uuid.UUID{}, err
		}

		publicKey := []byte(pk)

		decodedSignature, err := base64.StdEncoding.DecodeString(response.Signature)

		if err != nil {
			fmt.Println("Failed to decode Signature")
			fmt.Println(err)
			return uuid.UUID{}, err
		}

		dedcodedMessege, err := base64.StdEncoding.DecodeString(currentCode)

		if err != nil {
			fmt.Println("Failed to decode Signature")
			fmt.Println(err)
			return uuid.UUID{}, err

		}

		isValid := ed25519.Verify(publicKey, dedcodedMessege, decodedSignature)
		if isValid {
			mutated := authRequest
			mutated.Approved = true
			mutated.ApprovedKey = key
			authService.AuthRequests.CompareAndSwap(response.Uuid, authRequest, mutated)

			return authRequest.Id, nil
		}
	}

	return uuid.Nil, nil
}

func (authService *AuthenticationService) ExchangeCodeForToken(code uuid.UUID) (uuid.UUID, bool) {
	request, ok := authService.AuthRequests.Load(code)

	if !ok {
		return uuid.UUID{}, false
	}

	authRequest := request.(dtos.StoredAuthRequest)

	if !authRequest.Approved {
		return uuid.UUID{}, true
	}

	return authRequest.Id, true
}

func (authService *AuthenticationService) ExchangeCodeForPublicKey(code uuid.UUID) (string, bool) {
	request, ok := authService.AuthRequests.Load(code)

	if !ok {
		return "", false
	}

	authRequest := request.(dtos.StoredAuthRequest)

	if !authRequest.Approved {
		return "", true
	}

	return authRequest.ApprovedKey, true

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

	for range signal {
		fmt.Println("Something has happened at", time.Now())
		authService.AuthRequests.Range(func(key, value any) bool {
			storedRequest := value.(dtos.StoredAuthRequest)
			expired := storedRequest.Time.UTC().Before(time.Now().UTC())
			if expired {
				authService.AuthRequests.Delete(key)
			}

			return true
		})
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
