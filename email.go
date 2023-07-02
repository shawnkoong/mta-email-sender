package main

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender interface {
	SendEmail(subject string, body string, to []string, cc []string) error
}

type GmailSender struct {
	senderName     string
	senderAddress  string
	senderPassword string
}

func NewGmailSender(senderName string, senderAddress string, senderPassword string) EmailSender {
	return &GmailSender{
		senderName:     senderName,
		senderAddress:  senderAddress,
		senderPassword: senderPassword,
	}
}

func (gs *GmailSender) SendEmail(subject string, body string, to []string, cc []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", gs.senderName, gs.senderAddress)
	e.Subject = subject
	e.HTML = []byte(body)
	e.To = to
	e.Cc = cc

	auth := smtp.PlainAuth("", gs.senderAddress, gs.senderPassword, "smtp.gmail.com")
	return e.Send("smtp.gmail.com:587", auth)
}
