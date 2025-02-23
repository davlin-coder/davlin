package template

import (
	"embed"
)

//go:embed templates/verification_email.html
var verificationEmailTemplate embed.FS

// VerificationEmailData 验证码邮件模板数据
type VerificationEmailData struct {
	Code          string
	ExpireMinutes int
}

// InitDefaultTemplates 初始化默认模板
func (tm *templateManager) InitDefaultTemplates() error {
	templateContent, err := verificationEmailTemplate.ReadFile("templates/verification_email.html")
	if err != nil {
		return err
	}
	return tm.AddTemplate("verification_email", string(templateContent))
}
