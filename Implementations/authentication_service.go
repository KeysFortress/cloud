package implementations

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/syncmap"

	"leanmeal/api/dtos"
	"leanmeal/api/utils"
)

type AuthenticationService struct {
	AuthRequests syncmap.Map
	Domain       string
}

func (authService *AuthenticationService) GetMessage(email *string, id *uuid.UUID) dtos.InitAuthReponse {

	uuid := uuid.New()
	code := utils.GenerateRandomString(32)

	authResponse := dtos.StoredAuthRequest{
		Id:       *id,
		Name:     *email,
		Code:     code,
		Uuid:     uuid.String(),
		Time:     time.Now().UTC().Add(time.Duration(time.Minute * 10)),
		Approved: false,
		Ignore:   false,
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

func (authService *AuthenticationService) GetAuthRequest(id uuid.UUID) (dtos.StoredAuthRequest, error) {
	request, ok := authService.AuthRequests.Load(id)
	if !ok {
		return dtos.StoredAuthRequest{}, errors.New("request id missing")
	}

	authRequest := request.(dtos.StoredAuthRequest)
	return authRequest, nil
}

func (authService *AuthenticationService) StoreAuthRequest(request dtos.StoredAuthRequest) bool {
	authService.AuthRequests.Store(request.Id, request)
	return true
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
			if expired && !storedRequest.Ignore {
				authService.AuthRequests.Delete(key)
			}

			return true
		})
	}

}
