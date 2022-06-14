package model

import (
	"errors"
	"time"
)

const doneCompletion uint = 100

var ErrTodoCompletionTooHigh error = errors.New("cannot update, completion can not be greater than 100")
var ErrTodoCompletionMarkedAsDone error = errors.New("cannot update, todo already marked as done")

type Todo struct {
	Id          uint      `gorm:"primaryKey;auto increment;not null"`
	ExpiredAt   time.Time `gorm:"not null"`
	Title       string    `gorm:"not null;type:VARCHAR(255)"`
	Description string    `gorm:"not null;type:VARCHAR(255)"`
	Completion  uint
}

func NewTodo(title string, description string, expiredAt time.Time) *Todo {

	return &Todo{
		ExpiredAt:   expiredAt,
		Title:       title,
		Description: description,
		Completion:  0,
	}
}

func (t *Todo) MarkAsDone() error {
	if t.Completion == doneCompletion {
		return ErrTodoCompletionMarkedAsDone
	}
	t.Completion = doneCompletion
	return nil
}

func (t *Todo) Update(newExpiredAt time.Time, newTitle string, newDescription string, newCompletion uint) error {
	t.ExpiredAt = newExpiredAt
	t.Title = newTitle
	t.Description = newDescription
	if newCompletion > doneCompletion {
		return ErrTodoCompletionTooHigh
	}
	if t.Completion == doneCompletion {
		return ErrTodoCompletionMarkedAsDone
	}
	t.Completion = newCompletion
	return nil
}
