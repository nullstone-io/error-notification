package error_notification

import "context"

type User struct {
	Id       string
	Username string
	Email    string
	Token    interface{}
}

type userContextKey struct{}

func UserFromContext(ctx context.Context) *User {
	val, _ := ctx.Value(userContextKey{}).(*User)
	return val
}

func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey{}, user)
}
