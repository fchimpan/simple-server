package main

import (
	"context"
	"net/http"

	"github.com/fchimpan/simple-server/auth"
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

	clocker := clock.RealClocker{}
	v := validator.New()

	// init db
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}

	// init redis
	rcli, err := store.NewKVS(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}

	// init auth
	jwter, err := auth.NewJWTer(rcli, clocker)
	if err != nil {
		return nil, cleanup, err
	}

	// init repo
	repo := &store.Repository{Clocker: clocker}

	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: repo},
		Validator: v,
	}

	lt := &handler.ListTask{
		Service: &service.ListTasks{DB: db, Repo: repo},
	}

	ut := &handler.UpdateTask{
		Service:   &service.UpdateTask{DB: db, Repo: repo},
		Validator: v,
	}

	mux.Route("/tasks", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter))
		r.Get("/", lt.ServeHTTP)
		r.Post("/", at.ServeHTTP)
		r.Put("/{id}", ut.ServeHTTP)
	})

	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: repo},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)

	l := &handler.Login{
		Service: &service.Login{
			DB:             db,
			Repo:           repo,
			TokenGenerator: jwter,
		},
		Validator: v,
	}
	mux.Post("/login", l.ServeHTTP)

	mux.Route("/admin", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter), handler.AdminMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"message":"admin only"}`))
		})
	})

	return mux, nil, nil
}

// {"name":"user","password":"password","role":"admin"}
