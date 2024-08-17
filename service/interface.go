package service

import (
	"context"

	"github.com/fchimpan/simple-server/entity"
	"github.com/fchimpan/simple-server/store"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . TaskAdder TaskLister UserRegister UserGetter TokenGenerator TaskUpdater
type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer, id entity.UserID) (entity.Tasks, error)
}

type TaskUpdater interface {
	UpdateTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type UserRegister interface {
	RegisterUser(ctx context.Context, db store.Execer, u *entity.User) error
}

type UserGetter interface {
	GetUser(ctx context.Context, db store.Queryer, name string) (*entity.User, error)
}

type TokenGenerator interface {
	GenerateToken(ctx context.Context, u entity.User) ([]byte, error)
}
