package notifyService

import (
	"TodoList/internal/core/domain"
	"TodoList/internal/core/logger"
	"context"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type TemplateService interface {
	GetSubject(notifyType domain.NotifyType, data map[string]string) (string, error)
	GetContent(notifyType domain.NotifyType, data map[string]string) (string, error)
}

func (s *EmailSenderService) SendEmail(ctx context.Context, event domain.Notify) { //TODO переписать
	log := logger.FromContext(ctx)

	if event.Email == "" {
		log.Error("cannot send email: email is empty")
		return
	}

	subject, err := s.templateService.GetSubject(event.NotifyType, event.Data)
	if err != nil {
		log.Error("failed to build email subject", zap.Error(err))
		return
	}

	content, err := s.templateService.GetContent(event.NotifyType, event.Data)
	if err != nil {
		log.Error("failed to build email content", zap.Error(err))
		return
	}

	if err := s.sendMessage(event.Email, subject, content); err != nil {
		log.Error("failed to send email", zap.Error(err))
		return
	}

	log.Info("email sent successfully")
}

func (s *EmailSenderService) sendMessage(to, subject, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.fromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", content)

	return s.mailSender.DialAndSend(m)
}
