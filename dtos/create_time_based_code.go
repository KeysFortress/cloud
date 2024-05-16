package dtos

type CreateTimeBasedCode struct {
	Website   string
	Email     string
	Secret    string
	Type      int
	Validity  int
	Algorithm int
}
