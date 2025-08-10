package template

import (
	"os"
)

type SenderAdapter interface {
	SMSAdapter | MailAdapter
}

type SMSAdapter struct {
	Attachments []os.File
	To          string
	CC          []string
}

type MailAdapter struct {
	Attachments []os.File
	To          string
	From        string
	Subject     string
	CC          []string
}

type Sender[T SenderAdapter] interface {
	Send(text string) error
}
