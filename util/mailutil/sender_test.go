package mailutil

import (
	"crypto/tls"
	"github.com/karosown/katool-go/util/template"
	"testing"
)

func TestEmailClient_Send(t *testing.T) {
	client := EmailClient{
		MailAdapter: &template.MailAdapter{
			Attachments: nil,
			To:          "xxxxxx",
			From:        "xxxxxx",
			Subject:     "测试",
			CC:          nil,
		},
		EmailConfig: &EmailConfig{
			Identity: "",
			Username: "xxxxxx",
			Password: "xxxxxx",
			Host:     "smtp.qq.com",
			Port:     "465",
			TlsConfig: &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         "smtp.qq.com",
			},
		},
	}
	err := client.Send("<h1>测试，你好</h1>")
	if err != nil {
		t.Error(err)
	}
}
