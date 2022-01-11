package error_notification

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"runtime/debug"
	"time"
)

type notifierContextKey struct{}

func NotifierFromContext(ctx context.Context) *Notifier {
	val, _ := ctx.Value(notifierContextKey{}).(*Notifier)
	return val
}

func ContextWithNotifier(ctx context.Context, notifier *Notifier) context.Context {
	return context.WithValue(ctx, notifierContextKey{}, notifier)
}

type Notifier struct {
	Client *Client
}

type ResponseData interface {
	StatusCode() int
	Body() string
}

// NotifyHttpErrorHandler is designed to be able to be passed to Middleware.
// It will extract the user from the request, the error from the ResponseData, and the status_code from the ResponseData.
// It uses the default values for access token and environment.
func (n *Notifier) NotifyHttpErrorHandler(r *http.Request, data ResponseData, duration time.Duration) {
	if statusCode := data.StatusCode(); statusCode < 400 {
		return
	}
	vars := mux.Vars(r)
	user := UserFromContext(r.Context())
	var userToken interface{}
	if user != nil {
		userToken = user.Token
	}
	n.Client.NotifyError(user, data.Body(), map[string]interface{}{
		"api":            true,
		"request_method": r.Method,
		"request_uri":    r.URL,
		"org_name":       vars["orgName"],
		"request_id":     r.Header.Get("X-Request-Id"),
		"user_token":     userToken,
		"status_code":    data.StatusCode(),
		"duration":       duration,
	})
}

// NotifyHttpError will extract the user from the request, the error from the ResponseData, and the status_code from the ResponseData.
// It uses the default values for access token and environment.
func (n *Notifier) NotifyHttpError(r *http.Request, error interface{}) {
	vars := mux.Vars(r)
	user := UserFromContext(r.Context())
	var userToken interface{}
	if user != nil {
		userToken = user.Token
	}
	n.Client.NotifyError(user, fmt.Sprintf("%v", error), map[string]interface{}{
		"api":            true,
		"request_method": r.Method,
		"request_uri":    r.URL,
		"org_name":       vars["orgName"],
		"request_id":     r.Header.Get("X-Request-Id"),
		"user_token":     userToken,
	})
}

// NotifyHttpCriticalHandler is designed to be able to be passed to PanicMiddleware.
// It will extract the user from the request, the error from the ResponseData, and the status_code from the ResponseData.
// It uses the default values for access token and environment.
func (n *Notifier) NotifyHttpCriticalHandler(r *http.Request, error interface{}) {
	vars := mux.Vars(r)
	rawStack := debug.Stack()
	user := UserFromContext(r.Context())
	var userToken interface{}
	if user != nil {
		userToken = user.Token
	}
	n.Client.NotifyCritical(user, fmt.Sprintf("%v", error), map[string]interface{}{
		"api":            true,
		"request_method": r.Method,
		"request_uri":    r.URL,
		"org_name":       vars["orgName"],
		"request_id":     r.Header.Get("X-Request-Id"),
		"user_token":     userToken,
		"stack":          string(rawStack),
	})
}
