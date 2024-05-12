package interfaces

type JwtService interface {
	IssueToken(role string, id string, deviceKey string) map[string]any
	ValidateToken(token string) bool
	ExtractValue(token string, key string) interface{}
}
