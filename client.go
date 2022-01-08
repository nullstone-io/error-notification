package error_notification

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/rollbar/rollbar-go"
	"net/http"
	"os"
	"sync"
)

const (
	rollbarAccessTokenEnvVar = "ROLLBAR_ACCESS_TOKEN"
	nullstoneEnvEnvVar       = "NULLSTONE_ENV"
)

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

func AddErrorClientToRequest(client *Client) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r.WithContext(ContextWithErrorClient(r.Context(), client)))
		})
	}
}

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
