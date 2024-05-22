package interfaces

type JwtService interface {
	IssueMfaToken(id string, deviceKey string) map[string]any
	IssueToken(role string, id string, deviceKey string) map[string]any
	ValidateToken(token string) (bool, error)
	ValidateMfaToken(token string) (bool, error)
	ExtractValue(token string, key string) interface{}
}
