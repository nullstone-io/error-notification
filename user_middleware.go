package error_notification

import (
	"github.com/gorilla/mux"
	"net/http"
)

type GetUserFunc func(r *http.Request) *User

func SetUserMiddleware(fn GetUserFunc) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := fn(r)
			newRequest := r.WithContext(ContextWithUser(r.Context(), user))
			handler.ServeHTTP(w, newRequest)
		})
	}
}
