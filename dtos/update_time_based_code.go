package dtos

type UpdateTimeBasedCode struct {
	Id       string
	Website  string
	Email    string
	Secret   string
	Type     int
	Validity int
}
