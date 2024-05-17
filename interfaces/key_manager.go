package interfaces

type KeyManager interface {
	Generate(issuer string, size int) ([]byte, []byte, error)
	SignData(data []byte, privateKey []byte) ([]byte, error)
	VerifySignature(publicKey []byte, challange []byte, signature []byte) (bool, error)
}
