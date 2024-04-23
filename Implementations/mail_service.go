package implementations

import (
	"crypto/tls"
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

type MailService struct {
	From     string
	Password string
	Smtp     string
	Port     int
	SkipSSl  bool
}

func (mailService *MailService) SendMessage(to string, subject string, messege string) (bool, error) {

	m := gomail.NewMessage()
	m.SetHeader("From", mailService.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", messege)

	d := gomail.NewDialer(mailService.Smtp, mailService.Port, mailService.From, mailService.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: mailService.SkipSSl}

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Failed to send email to %v at %v ", to, time.Now().UTC())
		return false, err
	}

	return true, nil
}

func (mailService *MailService) SendTemplate(to string, subject string, template string) (bool, error) {

	m := gomail.NewMessage()
	m.SetHeader("From", mailService.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", template)

	d := gomail.NewDialer(mailService.Smtp, mailService.Port, mailService.From, mailService.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: mailService.SkipSSl}

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Failed to send email to %v at %v ", to, time.Now().UTC())
		return false, err
	}

	return true, nil
}
