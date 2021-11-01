package email

import (
	"crypto/tls"
	"wiki_bit/config"

	"gopkg.in/gomail.v2"
)

//邮件设置
type EmailMessage struct {
	From        string
	To          []string
	Cc          []string
	Subject     string
	ContentType string //text/plain text/html
	Content     string
	Attach      string
}

//邮箱client
type EmailClient struct {
	Host     string
	Port     int
	Username string
	Password string
	Message  *EmailMessage
}

//初始化邮箱
func NewEmailMessage(subject, contentType, content, attach string, to, cc []string) *EmailMessage {
	return &EmailMessage{
		From:        config.Conf().Email.UserName,
		Subject:     subject,
		ContentType: contentType,
		Content:     content,
		To:          to,
		Cc:          cc,
		Attach:      attach,
	}
}

//init new emailclient
func NewEmailClient(message *EmailMessage) *EmailClient {

	return &EmailClient{
		Host:     config.Conf().Email.Host,
		Port:     config.Conf().Email.Port,
		Username: config.Conf().Email.UserName,
		Password: config.Conf().Email.Password,
		Message:  message,
	}
}

//发送邮件
func (c *EmailClient) SendMessage() (bool, error) {
	e := gomail.NewPlainDialer(c.Host, c.Port, c.Username, c.Password)
	if 587 == c.Port {
		e.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	dm := gomail.NewMessage()
	dm.SetHeader("From", c.Message.From)
	dm.SetHeader("To", c.Message.To...)

	if len(c.Message.Cc) != 0 {
		dm.SetHeader("Cc", c.Message.Cc...)
	}

	dm.SetHeader("Subject", c.Message.Subject)
	dm.SetBody(c.Message.ContentType, c.Message.Content)

	if c.Message.Attach != "" {
		dm.Attach(c.Message.Attach)
	}

	if err := e.DialAndSend(dm); err != nil {
		return false, err
	}
	return true, nil
}
