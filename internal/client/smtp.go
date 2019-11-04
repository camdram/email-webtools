package client

import (
	"fmt"
	"log"
	"net/smtp"
)

type Mailer struct {
	Client *smtp.Client
}

func NewMailer() *Mailer {
	c, err := smtp.Dial("localhost:25")
	if err != nil {
		log.Fatal(err)
	}
	return &Mailer{
		Client: c,
	}
}

func (m *Mailer) Teardown() {
	err := m.Client.Quit()
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Mailer) Send(from string, to string, subject string, body string) {
	c := m.Client
	if err := c.Mail(from); err != nil {
		log.Fatal(err)
	}
	if err := c.Rcpt(to); err != nil {
		log.Fatal(err)
	}
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(wc, "Subject: [Camdram] "+subject+"\n")
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(wc, body)
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}
}
