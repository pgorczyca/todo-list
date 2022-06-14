package model_test

import (
	"testing"
	"time"

	"github.com/pgorczyca/todo-list/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTodo(t *testing.T) {
	t.Parallel()
	todoTitle := "test title"
	todoDescription := "test desc"
	todoExpiredAt := time.Now().Add(time.Hour * 48)

	todo := model.NewTodo(todoTitle, todoDescription, todoExpiredAt)

	assert.Equal(t, todoTitle, todo.Title)
	assert.Equal(t, todoDescription, todo.Description)
	assert.Equal(t, todoExpiredAt, todo.ExpiredAt)
	assert.Equal(t, uint(0), todo.Completion)
}
func TestMarkAsDone(t *testing.T) {
	t.Parallel()
	todo := model.NewTodo("test title", "test desc", time.Now().Add(time.Hour*48))

	err := todo.MarkAsDone()
	require.NoError(t, err)

	assert.Equal(t, uint(100), todo.Completion)
}

func TestMarkAsDoneThrowsErrorWhenAlreadyDone(t *testing.T) {
	t.Parallel()
	todo := model.NewTodo("test title", "test description", time.Now().Add(time.Hour*48))
	todo.MarkAsDone()

	err := todo.MarkAsDone()
	require.ErrorIs(t, err, model.ErrTodoCompletionMarkedAsDone)
}

func TestUpdateWithCorrectValues(t *testing.T) {
	t.Parallel()
	updatedExpiredAt := time.Now()
	updatedTitle := "updated title"
	updatedDescription := "updated desc"
	var updatedCompletion uint = 40
	todo := model.NewTodo("test title", "test desc", time.Now().Add(time.Hour*48))

	err := todo.Update(updatedExpiredAt, updatedTitle, updatedDescription, updatedCompletion)
	require.NoError(t, err)

	assert.Equal(t, updatedExpiredAt, todo.ExpiredAt)
	assert.Equal(t, updatedTitle, todo.Title)
	assert.Equal(t, updatedDescription, todo.Description)
	assert.Equal(t, updatedCompletion, todo.Completion)
}

func TestUpdateCompletionWithTooHighValueThrowsError(t *testing.T) {
	t.Parallel()
	updatedExpiredAt := time.Now()
	updatedTitle := "updated title"
	updatedDescription := "updated desc"
	var updatedCompletion uint = 400
	todo := model.NewTodo("test title", "test desc", time.Now().Add(time.Hour*48))

	err := todo.Update(updatedExpiredAt, updatedTitle, updatedDescription, updatedCompletion)
	require.ErrorIs(t, err, model.ErrTodoCompletionTooHigh)
}

func TestUpdateMarkAsDoneThrowsError(t *testing.T) {
	t.Parallel()
	todo := model.NewTodo("test title", "test desc", time.Now().Add(time.Hour*48))
	todo.MarkAsDone()

	err := todo.Update(time.Now(), "updated title", "updated desc", 10)
	require.ErrorIs(t, err, model.ErrTodoCompletionMarkedAsDone)
}
