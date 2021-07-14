package controllers

// SmtpService implements smtp related functions
type SmtpService interface {
	SendEmail(to []string, message []byte)
}
