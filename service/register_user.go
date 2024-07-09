package service

import (
	"context"
	"fmt"

	"github.com/fchimpan/simple-server/entity"
	"github.com/fchimpan/simple-server/store"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	DB   store.Execer
	Repo UserRegister
}

func (ru *RegisterUser) RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("fatal generate password: %w", err)
	}

	u := &entity.User{Name: name, Password: string(pw), Role: role}

	if err := ru.Repo.RegisterUser(ctx, ru.DB, u); err != nil {
		return nil, err
	}

	return u, nil
}
