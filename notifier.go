package error_notification

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"runtime/debug"
	"time"
)

type Notifier struct {
	GetUserFn      func(r *http.Request) *User
	GetUserTokenFn func(r *http.Request) string
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
	c := DefaultClient()
	c.NotifyError(n.GetUserFn(r), data.Body(), map[string]interface{}{
		"api":            true,
		"request_method": r.Method,
		"request_uri":    r.URL,
		"org_name":       vars["orgName"],
		"request_id":     r.Header.Get("X-Request-Id"),
		"user_token":     n.GetUserTokenFn(r),
		"status_code":    data.StatusCode(),
		"duration":       duration,
	})
}

// NotifyHttpError will extract the user from the request, the error from the ResponseData, and the status_code from the ResponseData.
// It uses the default values for access token and environment.
func (n *Notifier) NotifyHttpError(r *http.Request, error interface{}) {
	vars := mux.Vars(r)
	c := DefaultClient()
	c.NotifyError(n.GetUserFn(r), fmt.Sprintf("%v", error), map[string]interface{}{
		"api":            true,
		"request_method": r.Method,
		"request_uri":    r.URL,
		"org_name":       vars["orgName"],
		"request_id":     r.Header.Get("X-Request-Id"),
		"user_token":     n.GetUserTokenFn(r),
	})
}

// NotifyHttpCriticalHandler is designed to be able to be passed to PanicMiddleware.
// It will extract the user from the request, the error from the ResponseData, and the status_code from the ResponseData.
// It uses the default values for access token and environment.
func (n *Notifier) NotifyHttpCriticalHandler(r *http.Request, error interface{}) {
	vars := mux.Vars(r)
	c := DefaultClient()
	rawStack := debug.Stack()
	c.NotifyCritical(n.GetUserFn(r), fmt.Sprintf("%v", error), map[string]interface{}{
		"api":            true,
		"request_method": r.Method,
		"request_uri":    r.URL,
		"org_name":       vars["orgName"],
		"request_id":     r.Header.Get("X-Request-Id"),
		"user_token":     n.GetUserTokenFn(r),
		"stack":          string(rawStack),
	})
}
