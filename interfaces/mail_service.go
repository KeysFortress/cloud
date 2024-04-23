package interfaces

type MailService interface {
	SendMessage(to string, subject string, messege string) (bool, error)
	SendTemplate(to string, subject string, template string) (bool, error)
}
