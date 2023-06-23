package email

import (
	"encoding/json"
	"fmt"
	"mime/multipart"

	"gopkg.in/gomail.v2"
)

type Body struct {
	Emails  string                `form:"emails"`
	Subject string                `form:"subject"`
	Message string                `form:"message"`
	File    *multipart.FileHeader `form:"file"`
}

type Email struct {
	Name string `json:"email"`
}

// Send email with body data to one or more recipients
func Send(body Body) error {

	// Read JSON of recipients
	var emails []Email
	err := json.Unmarshal([]byte(body.Emails), &emails)
	if err != nil {
		return err
	}

	// Test email settings
	d := gomail.NewDialer("smtp.gmail.com", 587, "app.test607@gmail.com", "kpnbjiqainffceog")
	s, err := d.Dial()
	if err != nil {
		return err
	}

	// Send email to each recipient
	m := gomail.NewMessage()
	for _, r := range emails {
		m.SetHeader("From", "app.test607@gmail.com")
		m.SetAddressHeader("To", r.Name, r.Name)
		m.SetHeader("Subject", body.Subject)
		m.SetBody("text/html", fmt.Sprintf(body.Message))
		m.Attach("tmp/" + body.File.Filename)

		if err := gomail.Send(s, m); err != nil {
			return err
		}
		m.Reset()
	}

	return nil
}
