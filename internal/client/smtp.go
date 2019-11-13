package client

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	Client *smtp.Client
}

func NewMailer() (*Mailer, error) {
	c, err := smtp.Dial("localhost:25")
	if err != nil {
		return nil, err
	}
	return &Mailer{
		Client: c,
	}, nil
}

func (m *Mailer) Teardown() error {
	return m.Client.Quit()
}

func (m *Mailer) Send(from string, to string, subject string, body string) error {
	c := m.Client
	if err := c.Mail(from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}
	wc, err := c.Data()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(wc, "Subject: [Camdram] "+subject+"\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(wc, body)
	if err != nil {
		return err
	}
	return wc.Close()
}
