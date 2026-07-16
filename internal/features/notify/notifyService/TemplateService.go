package notifyService

import (
	"TodoList/internal/core/domain"
	"errors"
	"fmt"
)

var (
	errNoDataOrType = errors.New("notifyType or data is nil")
)

type TemplateServiceImp struct {
}

func (t *TemplateServiceImp) GetSubject(notifyType domain.NotifyType, data map[string]string) (string, error) {
	if notifyType == "" || data == nil {
		return "", errNoDataOrType
	}

	switch notifyType {
	case domain.NotifyTypeRegister:
		code := data["code"]
		if code == "" {
			code = ""
		}
		return fmt.Sprintf("TodoApp: Ваш код для входа [%s]", code), nil

	case domain.NotifyTypeReplayCode:
		code := data["code"]
		if code == "" {
			code = ""
		}
		return fmt.Sprintf("TodoApp: Повторный код [%s]", code), nil

	case domain.NotifyTypeLogin:
		return "TodoApp: Вход в аккаунт", nil

	default:
		return "Уведомление от TodoApp", nil
	}
}

func (t *TemplateServiceImp) GetContent(notifyType domain.NotifyType, data map[string]string) (string, error) {
	if notifyType == "" || data == nil {
		return "", errNoDataOrType
	}
	userName := data["userName"]
	if userName == "" {
		userName = "Name"
	}

	switch notifyType {
	case domain.NotifyTypeRegister:
		return fmt.Sprintf(`
Добро пожаловать в Kortex!

Введите этот код на странице подтверждения для завершения входа в ваш аккаунт.
Если вы не запрашивали вход, пожалуйста, проигнорируйте это письмо.

С уважением,Команда TodoApp`), nil

	case domain.NotifyTypeReplayCode:
		return fmt.Sprintf(`
Был запрошен повторный код

С уважением команда TodoApp`), nil

	case domain.NotifyTypeLogin:
		return fmt.Sprintf(`
Уважаемый %s,

В ваш аккаунт был выполнен вход.

Если это были не вы, пожалуйста, свяжитесь со службой поддержки.

С уважением,
Команда TodoApp`, userName), nil

	default:
		return "Уведомление от TodoApp", nil
	}
}
