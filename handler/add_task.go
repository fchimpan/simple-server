package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fchimpan/simple-server/entity"
	"github.com/fchimpan/simple-server/store"
	"github.com/go-playground/validator"
)

type AddTask struct {
	Store     *store.TaskStore
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Title string `json:"title" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	if err := at.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		},
			http.StatusInternalServerError,
		)
		return
	}

	t := &entity.Task{
		Title:     b.Title,
		Status:    entity.TaskStatusTodo,
		CreatedAt: time.Time{},
	}
	id, err := store.Tasks.Add(t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		},
			http.StatusInternalServerError,
		)
		return
	}

	resp := struct {
		ID entity.TaskID `json:"id"`
	}{
		ID: id,
	}
	RespondJSON(ctx, w, resp, http.StatusOK)

}
