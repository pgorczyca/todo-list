package app

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/pgorczyca/todo-list/internal/handler"
	"github.com/pgorczyca/todo-list/internal/model"
	"github.com/pgorczyca/todo-list/internal/repository"
	"github.com/pgorczyca/todo-list/internal/service"
	"github.com/pgorczyca/todo-list/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App interface {
	ServeHTTP()
}

type app struct {
	todoRepository repository.Todo
	todoService    service.Todo
	gormDB         *gorm.DB
	sqlDB          *sql.DB
}

var config = utils.GetConfig()

func New() App {
	utils.InitializeLogger()
	gormDB, err := gorm.Open(postgres.Open(config.PostgresDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := gormDB.DB()
	todoRepository := repository.NewTodoPostgres(gormDB)
	todoService := service.NewTodo(todoRepository)
	return &app{
		todoRepository: todoRepository,
		todoService:    todoService,
		gormDB:         gormDB,
		sqlDB:          sqlDB,
	}

}
func (a *app) ServeHTTP() {
	defer utils.Logger.Info("Stopping application")
	defer utils.Logger.Sync()
	defer a.sqlDB.Close()
	utils.Logger.Info("Running application")

	a.gormDB.AutoMigrate(&model.Todo{})

	router := gin.Default()
	todo := router.Group("/todos")
	{
		todo.POST("/", a.handleTodoService(handler.TodoCreate))
		todo.GET("/", a.handleTodoService(handler.TodoGetList))
		todo.GET("/:id", a.handleTodoService(handler.TodoGetByID))
		todo.PUT("/:id", a.handleTodoService(handler.TodoUpdate))
		todo.PATCH("/:id/done", a.handleTodoService(handler.TodoMarkAsDone))
		todo.DELETE("/:id", a.handleTodoService(handler.TodoDelete))

	}

	router.Run()
}

func (a *app) handleTodoService(handler func(c *gin.Context, s service.Todo)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c, a.todoService)
	}
}
