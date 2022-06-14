package service

import (
	"context"
	"time"

	"github.com/pgorczyca/todo-list/internal/dto"
	"github.com/pgorczyca/todo-list/internal/model"
	"github.com/pgorczyca/todo-list/internal/repository"
)

type Todo interface {
	Create(input dto.CreateTodoInput) (dto.TodoResponse, error)
	GetList(filter string) ([]dto.TodoResponse, error)
	GetByID(id uint) (dto.TodoResponse, error)
	MarkAsDone(id uint) (dto.TodoResponse, error)
	Update(input dto.UpdateTodoInput, id uint) (dto.TodoResponse, error)
	Delete(id uint) error
}

var ErrTodoCompletionMarkedAsDone error = model.ErrTodoCompletionMarkedAsDone
var ErrTodoCompletionTooHigh error = model.ErrTodoCompletionTooHigh

type todo struct {
	repository repository.Todo
}

func NewTodo(repo repository.Todo) Todo {
	return &todo{
		repository: repo,
	}
}

func (t *todo) Create(input dto.CreateTodoInput) (dto.TodoResponse, error) {
	model := model.NewTodo(input.Title, input.Description, input.ExpiredAt)
	if err := t.repository.Add(context.Background(), model); err != nil {
		return dto.TodoResponse{}, err
	}

	return mapModelToResponse(*model), nil
}

func (t *todo) GetList(filter string) ([]dto.TodoResponse, error) {
	var todos []*model.Todo
	var err error
	switch filter {
	case "today":
		y, m, d := time.Now().Date()
		today := time.Date(y, m, d, 0, 0, 0, 0, time.Now().Location())
		tomorrow := today.AddDate(0, 0, 1)
		todos, err = t.repository.GetIncoming(context.Background(), today, tomorrow)
	case "nextday":
		y, m, d := time.Now().Date()
		tomorrow := time.Date(y, m, d, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 1)
		nextday := tomorrow.AddDate(0, 0, 1)
		todos, err = t.repository.GetIncoming(context.Background(), tomorrow, nextday)
	case "currentweek":
		y, m, d := time.Now().Date()
		currentweek := time.Date(y, m, d, 0, 0, 0, 0, time.Now().Location())
		for currentweek.Weekday() != time.Monday {
			currentweek = currentweek.AddDate(0, 0, -1)
		}
		nextweek := currentweek.AddDate(0, 0, 7)
		todos, err = t.repository.GetIncoming(context.Background(), currentweek, nextweek)
	default:
		todos, err = t.repository.GetAll(context.Background())
	}

	if err != nil {
		return []dto.TodoResponse{}, err
	}
	todoResponses := []dto.TodoResponse{}
	for _, todo := range todos {
		response := mapModelToResponse(*todo)
		todoResponses = append(todoResponses, response)
	}

	return todoResponses, nil
}

func (t *todo) GetByID(id uint) (dto.TodoResponse, error) {
	todo, err := t.repository.GetByID(context.Background(), id)
	if err != nil {
		return dto.TodoResponse{}, err
	}
	return mapModelToResponse(*todo), nil
}

func (t *todo) MarkAsDone(id uint) (dto.TodoResponse, error) {
	todo, err := t.repository.GetByID(context.Background(), id)
	if err != nil {
		return dto.TodoResponse{}, err
	}
	if err = todo.MarkAsDone(); err != nil {
		return dto.TodoResponse{}, ErrTodoCompletionMarkedAsDone
	}
	if err = t.repository.Update(context.Background(), todo); err != nil {
		return dto.TodoResponse{}, err
	}
	return mapModelToResponse(*todo), nil
}
func (t *todo) Update(input dto.UpdateTodoInput, id uint) (dto.TodoResponse, error) {
	todo, err := t.repository.GetByID(context.Background(), id)
	if err != nil {
		return dto.TodoResponse{}, err
	}
	if err = todo.Update(input.ExpiredAt, input.Title, input.Description, input.Completion); err != nil {
		switch err {
		case ErrTodoCompletionTooHigh:
			return dto.TodoResponse{}, ErrTodoCompletionTooHigh
		case ErrTodoCompletionMarkedAsDone:
			return dto.TodoResponse{}, ErrTodoCompletionMarkedAsDone
		default:
			return dto.TodoResponse{}, err
		}
	}
	if err := t.repository.Update(context.Background(), todo); err != nil {
		return dto.TodoResponse{}, err
	}

	return mapModelToResponse(*todo), nil
}
func (t *todo) Delete(id uint) error {
	todo, err := t.repository.GetByID(context.Background(), id)
	if err != nil {
		return err
	}
	if err := t.repository.Delete(context.Background(), todo); err != nil {
		return err
	}
	return nil
}

func mapModelToResponse(model model.Todo) dto.TodoResponse {
	return dto.TodoResponse{
		Id:          model.Id,
		ExpiredAt:   model.ExpiredAt.Format(time.RFC3339),
		Title:       model.Title,
		Description: model.Description,
		Completion:  model.Completion,
		Done:        model.Completion == 100,
	}
}
