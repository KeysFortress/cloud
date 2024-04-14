package interfaces

type JwtService interface {
	IssueToken(role string, id string) string
	ValidateToken(token string) bool
	ExtractValue(token string, key string) interface{}
}
