package mailutil

import (
	"crypto/tls"
	"errors"
	"github.com/jordan-wright/email"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/util/template"
	"net/smtp"
	"net/textproto"
	"os"
)

type EmailConfig struct {
	Identity, Username, Password string
	Host, Port                   string
	TlsConfig                    *tls.Config
	FilterErrors                 []error
}
type EmailClient struct {
	*template.MailAdapter
	*EmailConfig
}

func (e *EmailConfig) makeServerAddr() string {
	return e.Host + ":" + e.Port
}
func (e *EmailConfig) makeAuth() smtp.Auth {
	return smtp.PlainAuth(e.Identity, e.Username, e.Password, e.Host)
}

func (e *EmailClient) Send(text string) error {
	if e.Host == "smtp.qq.com" || e.Host == "smtp.exmail.qq.com" {
		e.FilterErrors = append(e.FilterErrors, textproto.ProtocolError("short response: \u0000\u0000\u0000\u001A\u0000\u0000\u0000"))
		e.FilterErrors = stream.Of(&e.FilterErrors).Distinct().ToList()
	}
	newEmail := email.NewEmail()
	newEmail.From = e.From
	newEmail.To = append(e.CC, e.To)
	newEmail.Subject = e.Subject
	newEmail.Cc = e.CC
	newEmail.HTML = []byte(text)
	newEmail.Text = []byte(text)
	stream.Of(&e.Attachments).ForEach(func(item os.File) {
		_, err := newEmail.AttachFile(item.Name())
		if err != nil {
			return
		}
	})
	var err error
	if e.TlsConfig == nil {
		err = newEmail.Send(e.makeServerAddr(), e.makeAuth())
	} else {
		err = newEmail.SendWithTLS(e.makeServerAddr(), e.makeAuth(), e.TlsConfig)
	}
	if err != nil {
		if stream.Of(&e.FilterErrors).Filter(func(i error) bool {
			return errors.Is(err, i)
		}).Count() != 0 {
			return nil
		}
		return err
	}
	return nil
}
