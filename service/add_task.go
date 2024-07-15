package service

import (
	"context"
	"errors"

	"github.com/fchimpan/simple-server/auth"
	"github.com/fchimpan/simple-server/entity"
	"github.com/fchimpan/simple-server/store"
)

type AddTask struct {
	DB   store.Execer
	Repo TaskAdder
}

func (at *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	id, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, errors.New("failed to get user id")
	}

	t := &entity.Task{
		Title:  title,
		UserID: id,
		Status: entity.TaskStatusTodo,
	}
	err := at.Repo.AddTask(ctx, at.DB, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
