package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TaskID int64
type TaskStatus string

const (
	TaskStatusTodo    TaskStatus = "todo"
	TaskStatusDoing   TaskStatus = "doing"
	TaskStatusDone    TaskStatus = "done"
	TaskStatusWaiting TaskStatus = "waiting"
)

type Task struct {
	ID         TaskID     `json:"id" db:"id"`
	UserID     UserID     `json:"user_id" db:"user_id"`
	Title      string     `json:"title" db:"title"`
	Status     TaskStatus `json:"status" db:"status"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt time.Time  `json:"modified_at" db:"modified_at"`
}

type Tasks []*Task

type UserID int64
type User struct {
	ID         UserID    `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Password   string    `json:"password" db:"password"`
	Role       string    `json:"role" db:"role"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
}

func (u *User) ComparePassword(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
}
