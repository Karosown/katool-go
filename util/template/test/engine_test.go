package template

import (
	"crypto/tls"
	"github.com/karosown/katool-go/util/mailutil"
	"github.com/karosown/katool-go/util/randutil"
	"github.com/karosown/katool-go/util/template"
	"testing"
)

var hatmlayout = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <meta name="color-scheme" content="light dark">
  <meta name="supported-color-schemes" content="light dark">
  <title>Katool 验证码</title>
  <style>
    /* 限制宽度的容器在部分客户端需要 table 布局 + max-width hack */
    @media (prefers-color-scheme: dark) {
      body, .email-bg { background-color: #0b0c0f !important; }
      .email-container { background-color: #16181d !important; border-color: #2b2f36 !important; }
      .text, .muted { color: #e6e6e6 !important; }
      .muted { color: #a5acb8 !important; }
      .divider { border-color: #2b2f36 !important; }
    }
    /* 移动端小优化 */
    @media screen and (max-width: 600px) {
      .email-container { width: 100% !important; }
      .code { font-size: 28px !important; letter-spacing: 4px !important; }
      .btn { display: block !important; width: 100% !important; }
    }
  </style>
</head>
<body style="margin:0;padding:0;background:#f5f7fa;" class="email-bg">
  <!-- 预览文案（多数客户端会显示），通过隐藏样式避免正文显示 -->
  <div style="display:none;overflow:hidden;line-height:1;height:0;max-height:0;opacity:0;">
    您的验证码是 {{code}}，{{expiresInMinutes}} 分钟内有效，请勿泄露。
  </div>

  <table role="presentation" cellpadding="0" cellspacing="0" width="100%" style="background:#f5f7fa;" class="email-bg">
    <tr>
      <td align="center" style="padding:24px;">
        <!-- 容器 -->
        <table role="presentation" cellpadding="0" cellspacing="0" width="600" style="max-width:600px;width:100%;background:#ffffff;border:1px solid #e6e8eb;border-radius:10px;overflow:hidden;" class="email-container">
          <!-- 顶部 -->
          <tr>
            <td align="left" style="padding:20px 24px 0 24px;">
              <img src="{{logoUrl}}" width="120" height="auto" alt="Katool" style="display:block;border:none;outline:none;text-decoration:none;">
            </td>
          </tr>

          <!-- 标题 -->
          <tr>
            <td style="padding:12px 24px 0 24px;">
              <h1 style="margin:0;font-family:Arial,Helvetica,sans-serif;font-size:20px;line-height:28px;color:#111827;" class="text">
                验证您的邮箱
              </h1>
            </td>
          </tr>

          <!-- 提示 -->
          <tr>
            <td style="padding:8px 24px 0 24px;">
              <p style="margin:0;font-family:Arial,Helvetica,sans-serif;font-size:14px;line-height:22px;color:#374151;" class="text">
                您正在使用 Katool 进行邮箱验证。请输入以下验证码完成操作：
              </p>
            </td>
          </tr>

          <!-- 验证码块 -->
          <tr>
            <td align="center" style="padding:16px 24px 0 24px;">
              <div style="font-family:'Courier New',Consolas,monospace;font-weight:700;font-size:34px;letter-spacing:6px;color:#111827;background:#f2f4f7;border:1px dashed #d0d5dd;border-radius:8px;padding:16px 20px;" class="code">
                {{code}}
              </div>
            </td>
          </tr>

          <!-- 按钮 -->
          <tr>
            <td align="center" style="padding:16px 24px 0 24px;">
              <!-- 一些客户端不支持按钮，使用 a + 行内样式 -->
              <a href="{{actionUrl}}" target="_blank"
                 style="font-family:Arial,Helvetica,sans-serif;display:inline-block;background:#4a90e2;color:#ffffff;text-decoration:none;padding:12px 20px;border-radius:8px;font-size:14px;font-weight:600;"
                 class="btn">立即验证</a>
            </td>
          </tr>

          <!-- 说明 -->
          <tr>
            <td style="padding:14px 24px 0 24px;">
              <p style="margin:0;font-family:Arial,Helvetica,sans-serif;font-size:13px;line-height:20px;color:#6b7280;" class="muted">
                验证码将在 {{expiresInMinutes}} 分钟后失效。若非本人操作，请忽略此邮件。
              </p>
            </td>
          </tr>

          <!-- 备用链接 -->
          <tr>
            <td style="padding:14px 24px 0 24px;">
              <p style="margin:0;font-family:Arial,Helvetica,sans-serif;font-size:12px;line-height:20px;color:#6b7280;" class="muted">
                如果按钮无法点击，请将以下链接复制到浏览器打开：<br>
                <a href="{{actionUrl}}" style="color:#4a90e2;text-decoration:underline;word-break:break-all;">{{actionUrl}}</a>
              </p>
            </td>
          </tr>

          <!-- 分隔线 -->
          <tr>
            <td style="padding:16px 24px 0 24px;">
              <hr style="border:none;border-top:1px solid #e6e8eb;margin:0;" class="divider">
            </td>
          </tr>

          <!-- 支持与签名 -->
          <tr>
            <td style="padding:14px 24px 20px 24px;">
              <p style="margin:0;font-family:Arial,Helvetica,sans-serif;font-size:12px;line-height:20px;color:#6b7280;" class="muted">
                需要帮助？请联系 <a href="mailto:{{supportEmail}}" style="color:#4a90e2;text-decoration:underline;">{{supportEmail}}</a><br>
                请不要将验证码透露给任何人，包括自称“客服”的人员。
              </p>
            </td>
          </tr>

          <!-- 页脚 -->
          <tr>
            <td align="center" style="padding:14px 24px 24px 24px;">
              <p style="margin:0;font-family:Arial,Helvetica,sans-serif;font-size:12px;line-height:18px;color:#9aa1ab;" class="muted">
                © {{year}} Katool. All rights reserved.
              </p>
            </td>
          </tr>
        </table>

        <!-- 边距兼容（可选占位） -->
        <div style="height:24px; line-height:24px;">&#8202;</div>
      </td>
    </tr>
  </table>
</body>
</html>
`

func Test_SMS(t *testing.T) {
	engine := template.NewEngine[template.SMSAdapter](hatmlayout)
	engine.
		AddMappings(map[string]string{
			"code":             randutil.String(6),
			"expiresInMinutes": "5",
			"logoUrl":          "https://github.com/Karosown/katool-go/raw/main/logo.png",
			"actionUrl":        "https://github.com/Karosown/katool-go",
			"supportEmail":     "example@gmail.com",
			"year":             "2025",
		}).
		SetDelimiters("{{", "}}")
	err := engine.Send(&mailutil.EmailClient{
		MailAdapter: &template.MailAdapter{
			Attachments: nil,
			To:          "xxxxxx",
			From:        "xxxxxx",
			Subject:     "测试",
			CC:          nil,
		},
		EmailConfig: &mailutil.EmailConfig{
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
	})
	if err != nil {
		t.Error(err)
	}
}
