package dtos

type RequestPassword struct {
	UpperCase         bool
	LowerCase         bool
	Digits            bool
	Unique            bool
	SpecialCharacters bool
	Lenght            int
}
