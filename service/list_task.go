package service

import (
	"context"

	"github.com/fchimpan/simple-server/entity"
	"github.com/fchimpan/simple-server/store"
)

type ListTasks struct {
	DB   store.Queryer
	Repo TaskLister
}

func (lt *ListTasks) ListTasks(ctx context.Context) (entity.Tasks, error) {
	return lt.Repo.ListTasks(ctx, lt.DB)
}
