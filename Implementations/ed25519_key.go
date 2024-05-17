package implementations

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"errors"
	"fmt"
)

type ED25519Key struct {
}

func (e *ED25519Key) GenerateKeys(issuer string, size int) ([]byte, []byte, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Printf("Failed to generate key pair : %s", err)
		return []byte{}, []byte{}, err
	}

	privateBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	publicBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	return privateBytes, publicBytes, nil
}

func (e *ED25519Key) SignData(data []byte, privateKey []byte) ([]byte, error) {
	if len(privateKey) != ed25519.PrivateKeySize {
		return []byte{}, errors.New("invalid length for Ed25519 private key")
	}

	result := ed25519.Sign(privateKey, data)

	return result, nil
}

func (e *ED25519Key) VerifySignature(publicKey []byte, challange []byte, signature []byte) (bool, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return false, errors.New("invalid length for Ed25519 private key")
	}

	isValid := ed25519.Verify(publicKey, challange, signature)

	return isValid, nil
}
