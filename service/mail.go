package service

import (
	"os"
	"strings"

	"github.com/go-gomail/gomail"
)

// SendMail send mail according to the content
func SendMail(send string, to string, content string, sub string) error{
	// Get receivers and password
	toArray := strings.Split(to, ";")
	password := os.Getenv("ADMIN_MAIL_PASS")
	password = "T3Y2vAX3i1jH"
	// Set email content
	m := gomail.NewMessage()
	m.SetHeader("From", "admin@sysuactivity.com")
	m.SetHeader("To", toArray...)
	m.SetHeader("Subject", sub)
	m.SetBody("text/html", content)

	// Dial and send email
	d := gomail.NewDialer("smtp.exmail.qq.com", 25, "admin@sysuactivity.com", password)
	return d.DialAndSend(m)
}
