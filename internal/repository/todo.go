package repository

import (
	"context"
	"errors"
	"time"

	"github.com/pgorczyca/todo-list/internal/model"
	"github.com/pgorczyca/todo-list/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const timeout = time.Second * 5

type Todo interface {
	Add(ctx context.Context, todo *model.Todo) error
	GetAll(ctx context.Context) ([]*model.Todo, error)
	GetByID(ctx context.Context, id uint) (*model.Todo, error)
	GetIncoming(ctx context.Context, startTime time.Time, endTime time.Time) ([]*model.Todo, error)
	Update(ctx context.Context, model *model.Todo) error
	Delete(ctx context.Context, model *model.Todo) error
}

type todoPostgres struct {
	db *gorm.DB
}

type GormTodoModel struct {
	ID          uint      `gorm:"primaryKey;auto increment;not null"`
	ExpiredAt   time.Time `gorm:"not null"`
	Title       string    `gorm:"not null;type:VARCHAR(255)"`
	Description string    `gorm:"not null;type:VARCHAR(255)"`
	Completion  uint
}

func NewTodoPostgres(db *gorm.DB) Todo {
	return &todoPostgres{db: db}
}
func (p *todoPostgres) Add(ctx context.Context, todo *model.Todo) error {
	newctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	result := p.db.WithContext(newctx).Create(todo)
	if result.Error != nil {
		utils.Logger.Error("Not able to add.", zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func (p *todoPostgres) GetAll(ctx context.Context) ([]*model.Todo, error) {
	newctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	todos := []*model.Todo{}
	result := p.db.WithContext(newctx).Find(&todos)
	if result.Error != nil {
		utils.Logger.Error("Not able to get all.", zap.Error(result.Error))
		return nil, result.Error
	}
	return todos, nil
}

func (p *todoPostgres) GetByID(ctx context.Context, id uint) (*model.Todo, error) {
	newctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	todo := &model.Todo{}
	result := p.db.WithContext(newctx).Where("ID = ?", id).Find(todo)
	if result.RowsAffected == 0 {
		utils.Logger.Debug("Not able to get by id.", zap.Uint("id", id))
		return nil, errors.New("no records found")
	}
	return todo, nil
}

func (p *todoPostgres) GetIncoming(ctx context.Context, startTime time.Time, endTime time.Time) ([]*model.Todo, error) {
	newctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	todos := []*model.Todo{}
	result := p.db.WithContext(newctx).Where("expired_at BETWEEN ? AND ?", startTime, endTime).Find(&todos)
	if result.Error != nil {
		utils.Logger.Error("Not able to get incoming todos.", zap.Error(result.Error))
		return nil, result.Error
	}
	return todos, nil
}

func (p *todoPostgres) Update(ctx context.Context, model *model.Todo) error {
	newctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	saveresult := p.db.WithContext(newctx).Save(model)
	if saveresult.Error != nil {
		utils.Logger.Error("Not able to get save todos.", zap.Error(saveresult.Error))
		return saveresult.Error
	}
	return nil
}
func (p *todoPostgres) Delete(ctx context.Context, model *model.Todo) error {
	newctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result := p.db.WithContext(newctx).Delete(&model)
	if result.Error != nil {
		utils.Logger.Error("Not able to delete.", zap.Error(result.Error))
		return result.Error
	}
	return nil
}
