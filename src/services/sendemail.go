package services

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

// SendEmail ส่งอีเมลด้วย Gomail
// host: SMTP server host เช่น "smtp.gmail.com"
// port: SMTP port เช่น 587
// user: อีเมลผู้ส่ง
// pass: รหัสผ่านหรือ App Password
func SendEmail(host string, port int, to, subject, body string) error {
	
	user := os.Getenv("MAILER_EMAIL")
	pass := os.Getenv("MAILER_PASSWORD")

	if user == "" || pass == "" {
		return fmt.Errorf("email or password is empty")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body) // ถ้าต้องการ HTML ใช้ "text/html"

	d := gomail.NewDialer(host, port, user, pass)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
