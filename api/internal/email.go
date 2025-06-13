package notifystock

import (
	"context"

	"github.com/mailgun/mailgun-go/v5"
)

type MailGunClient struct {
	mg     *mailgun.Client
	domain string
}

var _ MailService = (*MailGunClient)(nil)

type MailGunClientConfig struct {
	Domain string
	ApiKey string
}

func NewMailGunClient(config MailGunClientConfig) *MailGunClient {
	mg := mailgun.NewMailgun(config.ApiKey)
	return &MailGunClient{
		domain: config.Domain,
		mg:     mg,
	}
}

func (m *MailGunClient) Send(from, to, subject, text string) error {
	message := mailgun.NewMessage(
		m.domain,
		from,
		subject,
		text,
		to,
	)
	res, err := m.mg.Send(context.Background(), message)
	if err != nil {
		return err
	}
	logger.Info("send email success", "status", res)
	return nil
}
