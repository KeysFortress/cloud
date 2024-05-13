package dtos

type IncomingPasswordRequest struct {
	Email             string
	Password          string
	Website           string
	UpperCase         bool
	LowerCase         bool
	Digits            bool
	Unique            bool
	SpecialCharacters bool
}
