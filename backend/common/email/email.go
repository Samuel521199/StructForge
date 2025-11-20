package email

import (
	"context"
	"fmt"
	"net/smtp"

	"StructForge/backend/common/log"
)

// Config 邮件配置
type Config struct {
	SMTPHost     string // SMTP 服务器地址
	SMTPPort     int    // SMTP 端口
	SMTPUser     string // SMTP 用户名
	SMTPPassword string // SMTP 密码
	FromEmail    string // 发件人邮箱
	FromName     string // 发件人名称
}

// EmailService 邮件服务接口
type EmailService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
	SendVerificationEmail(ctx context.Context, to, token string) error
	SendPasswordResetEmail(ctx context.Context, to, token string) error
}

// emailService 邮件服务实现
type emailService struct {
	config Config
}

// NewEmailService 创建邮件服务实例
func NewEmailService(config Config) EmailService {
	return &emailService{
		config: config,
	}
}

// SendEmail 发送邮件
func (s *emailService) SendEmail(ctx context.Context, to, subject, body string) error {
	// 构建邮件内容
	msg := []byte(fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("From: %s <%s>\r\n", s.config.FromName, s.config.FromEmail) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	// SMTP 认证
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)
	err := smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, msg)
	if err != nil {
		log.Error(ctx, "发送邮件失败",
			log.ErrorField(err),
			log.String("to", to),
			log.String("subject", subject),
		)
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	log.Info(ctx, "邮件发送成功",
		log.String("to", to),
		log.String("subject", subject),
	)

	return nil
}

// SendVerificationEmail 发送邮箱验证邮件
func (s *emailService) SendVerificationEmail(ctx context.Context, to, token string) error {
	// 构建验证链接（前端地址 + token）
	// TODO: 从配置中读取前端地址
	verifyURL := fmt.Sprintf("http://localhost:5173/auth/verify-email?token=%s", token)

	subject := "请验证您的邮箱 - StructForge"
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>邮箱验证</title>
		</head>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h1 style="color: #00FF00;">欢迎加入 StructForge</h1>
				<p>感谢您注册 StructForge 账号！</p>
				<p>请点击下面的链接验证您的邮箱地址：</p>
				<p style="text-align: center; margin: 30px 0;">
					<a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #00FF00; color: #000; text-decoration: none; border-radius: 4px; font-weight: bold;">验证邮箱</a>
				</p>
				<p>或者复制以下链接到浏览器中打开：</p>
				<p style="word-break: break-all; color: #666;">%s</p>
				<p style="color: #999; font-size: 12px; margin-top: 30px;">此链接将在 24 小时后过期。</p>
				<p style="color: #999; font-size: 12px;">如果您没有注册 StructForge 账号，请忽略此邮件。</p>
			</div>
		</body>
		</html>
	`, verifyURL, verifyURL)

	return s.SendEmail(ctx, to, subject, body)
}

// SendPasswordResetEmail 发送密码重置邮件
func (s *emailService) SendPasswordResetEmail(ctx context.Context, to, token string) error {
	// 构建重置链接（前端地址 + token）
	// TODO: 从配置中读取前端地址
	resetURL := fmt.Sprintf("http://localhost:5173/auth/reset-password?token=%s", token)

	subject := "重置您的密码 - StructForge"
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>重置密码</title>
		</head>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h1 style="color: #00FF00;">重置密码</h1>
				<p>您请求重置 StructForge 账号的密码。</p>
				<p>请点击下面的链接重置您的密码：</p>
				<p style="text-align: center; margin: 30px 0;">
					<a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #00FF00; color: #000; text-decoration: none; border-radius: 4px; font-weight: bold;">重置密码</a>
				</p>
				<p>或者复制以下链接到浏览器中打开：</p>
				<p style="word-break: break-all; color: #666;">%s</p>
				<p style="color: #999; font-size: 12px; margin-top: 30px;">此链接将在 1 小时后过期。</p>
				<p style="color: #999; font-size: 12px;">如果您没有请求重置密码，请忽略此邮件，您的密码将不会被更改。</p>
			</div>
		</body>
		</html>
	`, resetURL, resetURL)

	return s.SendEmail(ctx, to, subject, body)
}

// DefaultConfig 返回默认邮件配置
func DefaultConfig() Config {
	return Config{
		SMTPHost:     "smtp.gmail.com",
		SMTPPort:     587,
		SMTPUser:     "",
		SMTPPassword: "",
		FromEmail:    "noreply@structforge.com",
		FromName:     "StructForge",
	}
}
