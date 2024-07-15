package fixture

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/fchimpan/simple-server/entity"
)

func User(u *entity.User) *entity.User {
	res := &entity.User{
		ID:         entity.UserID(rand.Int()),
		Name:       "test" + strconv.Itoa(rand.Int())[0:5],
		Password:   "password",
		Role:       "admin",
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	if u == nil {
		return res
	}
	if u.ID != 0 {
		res.ID = u.ID
	}
	if u.Name != "" {
		res.Name = u.Name
	}
	if u.Password != "" {
		res.Password = u.Password
	}
	if u.Role != "" {
		res.Role = u.Role
	}
	if !u.CreatedAt.IsZero() {
		res.CreatedAt = u.CreatedAt
	}
	if !u.ModifiedAt.IsZero() {
		res.ModifiedAt = u.ModifiedAt
	}
	return res
}
