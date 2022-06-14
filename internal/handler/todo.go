package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pgorczyca/todo-list/internal/dto"
	"github.com/pgorczyca/todo-list/internal/service"
	"github.com/pgorczyca/todo-list/internal/utils"
	"github.com/pgorczyca/todo-list/internal/validation"
	"go.uber.org/zap"
)

func TodoCreate(c *gin.Context, s service.Todo) {
	var input dto.CreateTodoInput

	if err := c.ShouldBind(&input); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			c.JSON(http.StatusBadRequest, gin.H{"errors:": validation.Format(verr)})

		}

		utils.Logger.Debug("Not able to unmarshal.", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error:": "not able to unmarshal"})
		return
	}

	todoResponse, err := s.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "not able to create todo"})
		return
	}

	c.JSON(http.StatusCreated, todoResponse)
}

func TodoGetList(c *gin.Context, s service.Todo) {
	filter, _ := c.GetQuery("filter")
	todoResponse, err := s.GetList(filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "not able to get todo list"})
	}
	c.JSON(http.StatusOK, todoResponse)
}

func TodoGetByID(c *gin.Context, s service.Todo) {
	id, err := parseUintFromIdParam(c)
	if err != nil {
		return
	}
	response, err := s.GetByID(id)
	if err != nil {
		utils.Logger.Debug("Not able to find todo with given id.", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error:": "not able to find todo with given id"})
		return
	}
	c.JSON(http.StatusOK, response)
}

func TodoUpdate(c *gin.Context, s service.Todo) {
	id, err := parseUintFromIdParam(c)
	if err != nil {
		return
	}

	var input dto.UpdateTodoInput
	if err := c.ShouldBind(&input); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			utils.Logger.Debug("Not able to validate input struct.", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"errors:": validation.Format(verr)})
			return
		}
		utils.Logger.Debug("Not able to unmarshal.", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error:": "not able to unmarshal"})
		return
	}

	response, err := s.Update(input, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "not able to update todo"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func TodoMarkAsDone(c *gin.Context, s service.Todo) {
	id, err := parseUintFromIdParam(c)

	if err != nil {
		return
	}
	response, err := s.MarkAsDone(id)

	if err != nil {
		switch err {
		case service.ErrTodoCompletionMarkedAsDone:
			c.JSON(http.StatusBadRequest, gin.H{"error:": "todo is already marked as done"})
			return
		case service.ErrTodoCompletionTooHigh:
			c.JSON(http.StatusBadRequest, gin.H{"error:": "todo completion cannot be higher than 100"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "not able to mark todo as done"})
		return
	}

	c.JSON(http.StatusOK, response)

}

func TodoDelete(c *gin.Context, s service.Todo) {
	id, err := parseUintFromIdParam(c)
	if err != nil {
		return
	}
	if err = s.Delete(id); err != nil {
		utils.Logger.Debug("Not able to delete todo.", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error:": "not able to delete todo"})
		return
	}
	c.JSON(http.StatusOK, "todo was deleted.")

}

func parseUintFromIdParam(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Debug("Not able to parse id parameter.", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error:": "not able to parse id parameter"})
		return 0, errors.New("not able to parse id parameter")

	}
	return uint(id), nil
}
