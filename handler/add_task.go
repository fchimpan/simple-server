package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fchimpan/simple-server/entity"
	"github.com/go-playground/validator"
)

type AddTask struct {
	Service   AddTaskService
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

	t, err := at.Service.AddTask(ctx, b.Title)
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
		ID: t.ID,
	}
	RespondJSON(ctx, w, resp, http.StatusOK)

}
