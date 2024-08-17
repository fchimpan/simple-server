package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fchimpan/simple-server/entity"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type UpdateTask struct {
	Service   UpdateTaskService
	Validator *validator.Validate
}

func (ut *UpdateTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		ID     entity.TaskID     `json:"id" validate:"required"`
		Title  string            `json:"title"`
		Status entity.TaskStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	taskID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		},

			http.StatusInternalServerError,
		)
		return
	}
	b.ID = entity.TaskID(taskID)

	if err := ut.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		},
			http.StatusInternalServerError,
		)
		return
	}

	t, err := ut.Service.UpdateTask(ctx, b.ID, b.Title, b.Status)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
			Details: "failed to update task",
		},
			http.StatusInternalServerError,
		)
		return
	}

	resp := struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Status string `json:"status"`
	}{
		ID:     strconv.Itoa(int(t.ID)),
		Title:  t.Title,
		Status: string(t.Status),
	}
	RespondJSON(ctx, w, resp, http.StatusOK)

}
