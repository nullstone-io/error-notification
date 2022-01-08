package error_notification

import (
	"context"
	"github.com/rollbar/rollbar-go"
	"os"
)

const (
	rollbarAccessTokenEnvVar = "ROLLBAR_ACCESS_TOKEN"
	nullstoneEnvEnvVar       = "NULLSTONE_ENV"
)

type Client struct {
	AccessToken string
	Environment string
}

func DefaultClient() *Client {
	return &Client{
		AccessToken: os.Getenv(rollbarAccessTokenEnvVar),
		Environment: os.Getenv(nullstoneEnvEnvVar),
	}
}

func (c *Client) newRollbarClient() *rollbar.Client {
	var hostname, _ = os.Hostname()
	return rollbar.New(c.AccessToken, c.Environment, "", hostname, "")
}

func (c *Client) NotifyError(user *User, error string, extras map[string]interface{}) {
	rb := c.newRollbarClient()
	if user != nil {
		rb.SetPerson(user.Id, user.Username, user.Email)
	}
	rb.MessageWithExtras(rollbar.ERR, error, extras)
}

func (c *Client) NotifyCritical(user *User, error string, extras map[string]interface{}) {
	rb := c.newRollbarClient()
	if user != nil {
		rb.SetPerson(user.Id, user.Username, user.Email)
	}
	rb.MessageWithExtras(rollbar.CRIT, error, extras)
}

type clientContextKey struct{}

func ContextWithErrorClient(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, clientContextKey{}, client)
}

func ErrorClientFromContext(ctx context.Context) *Client {
	if val, ok := ctx.Value(clientContextKey{}).(*Client); ok {
		return val
	}
	return nil
}
