package controllers

type SmtpService interface{
	SendEmail(to []string, message []byte)
}
