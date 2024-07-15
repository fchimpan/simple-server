package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type Login struct {
	Service   LoginService
	Validator *validator.Validate
}

func (l *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		Username string `json:"user_name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		},
			http.StatusInternalServerError,
		)
		return
	}
	err := l.Validator.Struct(body)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		},
			http.StatusBadRequest,
		)
		return
	}

	jwt, err := l.Service.Login(ctx, body.Username, body.Password)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		},
			http.StatusInternalServerError,
		)
		return
	}
	resp := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: jwt,
	}

	RespondJSON(r.Context(), w, resp, http.StatusOK)
}
