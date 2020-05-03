package email

import (
	"fmt"
	"gopkg.in/mailgun/mailgun-go.v1"
	"net/url"
)

const (
	welcomeSubject = "Wubbalubbadubdub!"
	resetBaseURL   = "https://jarvis.ge/reset"
)

const welcomeText = `Hit the sack, Jack!

And that's why I always say, 'Shumshumschilpiddydah!

Rikitikitavi, bitch!
`

const resetTextTmpl = `Mr. Meeseeeeeeeeks

Please click link below:

%s

Token:

%s

Best,
God!
`

type Client struct {
	from string
	mg   mailgun.Mailgun
}

type ClientConfig func(*Client)

func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		from: "Jarvis Ge <demo@sandboxe83804ac69bf46a88f406b0194c11ed8.mailgun.org>",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

func WithMailgun(domain, apiKey, publicKey string) ClientConfig {
	return func(c *Client) {
		mg := mailgun.NewMailgun(domain, apiKey, publicKey)
		c.mg = mg
	}
}

func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.from = buildEmail(name, email)
	}
}

func (c *Client) Welcome(toName, toEmail string) error {
	message := mailgun.NewMessage(c.from, welcomeSubject, welcomeText, buildEmail(toName, toEmail))
	_, _, err := c.mg.Send(message)
	return err
}

func (c *Client) ResetPw(toEmail, token string) error {
	v := url.Values{}
	v.Set("token", token)
	resetUrl := resetBaseURL + "?" + v.Encode()
	resetText := fmt.Sprintf(resetTextTmpl, resetUrl, token)
	message := mailgun.NewMessage(c.from, welcomeSubject, resetText, toEmail)
	_, _, err := c.mg.Send(message)
	return err
}

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
