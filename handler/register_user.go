package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fchimpan/simple-server/entity"
	"github.com/go-playground/validator"
)

type RegisterUser struct {
	Service   RegisterUserService
	Validator *validator.Validate
}

func (ru *RegisterUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	if err := ru.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		},
			http.StatusInternalServerError,
		)
		return
	}

	u, err := ru.Service.RegisterUser(ctx, b.Name, b.Password, b.Role)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		},
			http.StatusInternalServerError,
		)
		return
	}
	resp := struct {
		ID entity.UserID `json:"id"`
	}{ID: u.ID}

	RespondJSON(ctx, w, resp, http.StatusOK)
}
