package notifystock

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type From struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
type To struct {
	Email string `json:"email"`
}

type MailTrap struct {
	From     From   `json:"from"`
	To       []To   `json:"to"`
	Subject  string `json:"subject"`
	Text     string `json:"text"`
	Category string `json:"category"`
}

func NewMailTrap(from, to, subject, text string) MailTrap {
	return MailTrap{
		From: From{
			Email: from,
			Name:  "MailTrap",
		},
		To: []To{
			{Email: to},
		},
		Subject:  subject,
		Text:     text,
		Category: "Integration Test",
	}
}

type MailTrapClient struct {
	url   string
	token string
}

func NewMailTrapClient(token string) *MailTrapClient {
	if !strings.HasPrefix(token, "Bearer") {
		token = "Bearer " + token
	}
	return &MailTrapClient{
		url:   "https://send.api.mailtrap.io/api/send",
		token: token,
	}
}

func (m *MailTrapClient) Send(from, to, subject, text string) error {
	mail := NewMailTrap(from, to, subject, text)
	payload, err := json.Marshal(mail)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", m.url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", m.token)
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	logger.Info("send email", "body", string(body))
	return nil
}
