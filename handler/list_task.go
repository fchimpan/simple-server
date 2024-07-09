package handler

import (
	"net/http"

	"github.com/fchimpan/simple-server/entity"
)

type ListTask struct {
	Service ListTasksService
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tasks, err := lt.Service.ListTasks(r.Context())
	if err != nil {
		RespondJSON(r.Context(), w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	resp := []task{}
	for _, t := range tasks {
		resp = append(resp, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
	}
	RespondJSON(r.Context(), w, resp, http.StatusOK)
}
