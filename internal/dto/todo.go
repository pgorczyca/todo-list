package dto

import "time"

type CreateTodoInput struct {
	ExpiredAt   time.Time `json:"expired_at" binding:"required"`
	Title       string    `json:"title" binding:"required,ascii,min=1,max=255"`
	Description string    `json:"description" binding:"required,min=1,max=255"`
}

type UpdateTodoInput struct {
	CreateTodoInput
	Completion uint `json:"completion" binding:"required,min=0,max=100"`
}

type TodoResponse struct {
	Id          uint   `json:"id"`
	ExpiredAt   string `json:"expired_at"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completion  uint   `json:"completion"`
	Done        bool   `json:"done"`
}
