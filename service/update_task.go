package service

import (
	"context"
	"errors"

	"github.com/fchimpan/simple-server/auth"
	"github.com/fchimpan/simple-server/entity"
	"github.com/fchimpan/simple-server/store"
)

type UpdateTask struct {
	DB   store.Execer
	Repo TaskUpdater
}

func (ut *UpdateTask) UpdateTask(ctx context.Context, id entity.TaskID, title string, status entity.TaskStatus) (*entity.Task, error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, errors.New("failed to get user id")
	}

	if err := ut.Repo.UpdateTask(ctx, ut.DB, &entity.Task{
		ID:     id,
		Title:  title,
		Status: status,
		UserID: userID,
	}); err != nil {
		return nil, err
	}

	return &entity.Task{
		ID:     id,
		Title:  title,
		Status: status,
	}, nil

}
