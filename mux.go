package main

import (
	"context"
	"net/http"

	"github.com/fchimpan/simple-server/clock"
	"github.com/fchimpan/simple-server/config"
	"github.com/fchimpan/simple-server/handler"
	"github.com/fchimpan/simple-server/service"
	"github.com/fchimpan/simple-server/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	repo := &store.Repository{Clocker: clock.RealClocker{}}

	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: repo},
		Validator: v,
	}
	mux.Post("/tasks", at.ServeHTTP)

	lt := &handler.ListTask{
		Service: &service.ListTasks{DB: db, Repo: repo},
	}
	mux.Get("/tasks", lt.ServeHTTP)

	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: repo},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)

	return mux, nil, nil
}

// {"name":"user","password":"password","role":"admin"}
