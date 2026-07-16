package notifyService

import (
	"gopkg.in/gomail.v2"
)

type MailSender interface {
	DialAndSend(m ...*gomail.Message) error
}

type EmailSenderService struct {
	mailSender      MailSender
	templateService TemplateService
	fromEmail       string
}

func NewEmailSenderService(mailSender MailSender, templateService TemplateService, fromEmail string) *EmailSenderService {
	return &EmailSenderService{
		mailSender:      mailSender,
		templateService: templateService,
		fromEmail:       fromEmail,
	}
}
