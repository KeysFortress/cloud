package implementations

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"fmt"
)

type RSAKey struct {
}

func (r *RSAKey) GenerateKeys(issuer string, size int) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		fmt.Printf("Failed to generate RSA key pair: %s\n", err)
		return []byte{}, []byte{}, err
	}

	privateBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	publicKey := &privateKey.PublicKey
	publicBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	return privateBytes, publicBytes, nil
}

func (r *RSAKey) SignData(data []byte, privateKey []byte) ([]byte, error) {
	priv, err := x509.ParsePKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	rsaPriv, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}

	hashed := sha256.Sum256(data)

	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPriv, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func (r *RSAKey) VerifySignature(publicKey []byte, challenge []byte, signature []byte) (bool, error) {
	pub, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return false, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("not an RSA public key")
	}

	hashed := sha256.Sum256(challenge)

	err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return false, err
	}

	return true, nil
}
