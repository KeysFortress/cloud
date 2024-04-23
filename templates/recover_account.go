package EmailTemplates

import (
	"bytes"
	"html/template"
	"log"
)

type RecoverAccount struct {
	Email string
	Code  string
}

func (ra *RecoverAccount) InitTemplate() string {
	t := template.New("email_templates/recover_account.html")

	var err error
	t, err = t.ParseFiles("email_templates/recover_account.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, &ra); err != nil {
		log.Println(err)
	}

	result := tpl.String()
	return result
}
