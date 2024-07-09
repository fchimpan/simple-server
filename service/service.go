package service

import (
	"context"

	"github.com/fchimpan/simple-server/entity"
	"github.com/fchimpan/simple-server/store"
)

type AddTask struct {
	DB   store.Execer
	Repo TaskAdder
}

func (at *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	t := &entity.Task{Title: title, Status: entity.TaskStatusTodo}
	err := at.Repo.AddTask(ctx, at.DB, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
