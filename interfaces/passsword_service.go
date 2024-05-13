package interfaces

type PasswordService interface {
	GeneratePassword(length int, lowerCase, upperCase, unique, specialChars bool) (string, error)
}
