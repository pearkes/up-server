package main

// Includes the logic to send notification emails.

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"time"
)

const crlf = "\r\n"

type Mail struct {
	Recipient string
	Subject   string
	Body      string
	From      string
	Content   []byte
}

type MailAuth struct {
	Sender   string
	Password string
	Host     string
	Auth     smtp.Auth
}

func getEmailAuth(auth MailAuth) smtp.Auth {
	au := smtp.PlainAuth(
		"",
		auth.Sender,
		auth.Password,
		auth.Host,
	)
	return au
}

func formatEmailBody(mail Mail) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "From: %s%s", mail.From, crlf)
	fmt.Fprintf(&buf, "To: %s%s", mail.Recipient, crlf)
	fmt.Fprintf(&buf, "Date: %s%s", time.Now().UTC().Format(time.RFC822), crlf)
	fmt.Fprintf(&buf, "Subject: %s%s%s", mail.Subject, crlf, mail.Body)
	return []byte(buf.String())
}

func sendMail(mail Mail) {
	mail.Content = formatEmailBody(mail)
	auth := getEmailConfig()
	auth.Auth = getEmailAuth(auth)
	if getDebug() == true {
		fmt.Println("Sending debug mail:\n\n", mail.Subject, mail.Recipient)
		return
	}
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		auth.Host+":25",
		auth.Auth,
		auth.Sender,
		[]string{mail.Recipient},
		mail.Content,
	)
	if err != nil {
		log.Fatal(err)
	}
}
