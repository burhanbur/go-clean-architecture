package utils

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

type IEmailRepository interface {
	SendEmail(email, otp string) error
}

type emailRepository struct {
}

func (e emailRepository) SendEmail(email string, otp string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", os.Getenv("MAIL_FROM_NAME"))
	mailer.SetHeader("To", email)
	mailer.SetAddressHeader("Cc", "tanisample@gmail.com", "Tani Sample")
	mailer.SetHeader("Subject", "OTP Login to Tani")
	mailer.SetBody("text/html", fmt.Sprintf("Hello,this is your otp number : <b>%v</b>", otp))

	port, _ := getenvInt(os.Getenv("MAIL_PORT"))

	dialer := gomail.NewDialer(
		os.Getenv("MAIL_HOST"),
		port,
		os.Getenv("MAIL_USERNAME"),
		os.Getenv("MAIL_PASSWORD"),
	)

	err := dialer.DialAndSend(mailer)

	if err != nil {
		Log{}.Error(err.Error())
	}

	return nil
}

func NewEmailRepository() IEmailRepository {
	return &emailRepository{}
}
