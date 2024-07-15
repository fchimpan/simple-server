package auth

import "context"

type LoginService interface {
	Login(ctx context.Context, username, password string) (string, error)
}
