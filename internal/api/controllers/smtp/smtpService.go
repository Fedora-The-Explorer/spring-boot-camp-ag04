package smtp

import (
	"fmt"
	"net/smtp"
)

const smtpHost = "email-smtp.eu-west-1.amazonaws.com"
const smtpPort = "587"
const smtpPassword = "BDVKBQLtJH5DFJ3isJMP80afrFXxjyIOKlMNrdyHw7aD"
const smtpUsername = "AKIA3QRJDSTT4P7LI2NJ"
const from = "luka.curic@ag04.io"


type SmtpService struct {
}
// NewEmailService creates a new instance of HeistValidator
func NewEmailService () *SmtpService{
	return &SmtpService{
	}
}

func (s *SmtpService) SendEmail(to []string, message []byte){

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email sent successfully!")
}