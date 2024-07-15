package service

import (
	"context"
	"errors"

	"github.com/fchimpan/simple-server/auth"
	"github.com/fchimpan/simple-server/entity"
	"github.com/fchimpan/simple-server/store"
)

type ListTasks struct {
	DB   store.Queryer
	Repo TaskLister
}

func (lt *ListTasks) ListTasks(ctx context.Context) (entity.Tasks, error) {
	id, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, errors.New("failed to get user id")
	}
	return lt.Repo.ListTasks(ctx, lt.DB, id)
}
