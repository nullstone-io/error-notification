package error_notification

import (
	"github.com/rollbar/rollbar-go"
	"os"
	"sync"
)

const (
	rollbarAccessTokenEnvVar = "ROLLBAR_ACCESS_TOKEN"
	nullstoneEnvEnvVar       = "NULLSTONE_ENV"
)

type Client struct {
	AccessToken string
	Environment string
	sync.Once
	rollbarClient *rollbar.Client
}

func DefaultClient() *Client {
	return &Client{
		AccessToken: os.Getenv(rollbarAccessTokenEnvVar),
		Environment: os.Getenv(nullstoneEnvEnvVar),
	}
}

func (c *Client) getRollbarClient() *rollbar.Client {
	c.Do(func() {
		var hostname, _ = os.Hostname()
		c.rollbarClient = rollbar.New(c.AccessToken, c.Environment, "", hostname, "")
	})
	return c.rollbarClient
}

func (c *Client) NotifyError(user *User, error string, extras map[string]interface{}) {
	rb := c.getRollbarClient()
	if user != nil {
		rb.SetPerson(user.Id, user.Username, user.Email)
	}
	rb.MessageWithExtras(rollbar.ERR, error, extras)
}

func (c *Client) NotifyCritical(user *User, error string, extras map[string]interface{}) {
	rb := c.getRollbarClient()
	if user != nil {
		rb.SetPerson(user.Id, user.Username, user.Email)
	}
	rb.MessageWithExtras(rollbar.CRIT, error, extras)
}
