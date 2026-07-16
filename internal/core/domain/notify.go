package domain

type NotifyType string

const (
	NotifyTypeRegister   NotifyType = "register"
	NotifyTypeReplayCode NotifyType = "replayCode"
	NotifyTypeLogin      NotifyType = "signIn"
)

type Notify struct {
	Email      string
	NotifyType NotifyType
	Data       map[string]string
}

func NewNotify(email string, notifyType NotifyType, data map[string]string) *Notify {
	return &Notify{
		Email:      email,
		NotifyType: notifyType,
		Data:       data,
	}
}
