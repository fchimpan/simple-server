package handler

import (
	"net/http"

	"github.com/fchimpan/simple-server/entity"
	"github.com/fchimpan/simple-server/store"
)

type ListTask struct {
	Store *store.TaskStore
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tasks := lt.Store.All()
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
