package interfaces

import (
	"github.com/google/uuid"

	"leanmeal/api/dtos"
)

type AuthenticationService interface {
	Start()
	GetMessage(email *string, id *uuid.UUID) dtos.InitAuthReponse
	VerifySignature(response dtos.FinishAuthResponse, keys *[]string) (uuid.UUID, error)
	ExchangeCodeForToken(code uuid.UUID) (uuid.UUID, bool)
	GetRequestById(id uuid.UUID) (uuid.UUID, error)
	ExchangeCodeForPublicKey(code uuid.UUID) (string, bool)
	StoreAuthRequest(request dtos.StoredAuthRequest) bool
	GetAuthRequest(id uuid.UUID) (dtos.StoredAuthRequest, error)
}
