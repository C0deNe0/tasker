package service

import (
	"context"

	"github.com/C0deNe0/go-tasker/internal/errs"
	"github.com/C0deNe0/go-tasker/internal/middleware"
	"github.com/C0deNe0/go-tasker/internal/model/todo"
	"github.com/C0deNe0/go-tasker/internal/repository"
	"github.com/C0deNe0/go-tasker/internal/server"
	"github.com/labstack/echo/v4"
)

type TodoService struct {
	server       *server.Server
	todoRepo     repository.TodoRepository
	categoryRepo repository.CategoryRepository
}

func NewTodoService(server *server.Server, todoRepo *repository.TodoRepository, categroyRepo *repository.CategoryRepository) *TodoService {
	return &TodoService{
		server:       server,
		todoRepo:     *todoRepo,
		categoryRepo: *categroyRepo,
	}
}

func (s *TodoService) CreateTodo(ctx echo.Context, userID string, payload *todo.CreateTodoPayload) (*todo.Todo, error) {
	logger := middleware.GetLogger(ctx)

	if payload.ParentTodoID != nil {
		parentTodo, err := s.todoRepo.CheckTodoExists(ctx.Request().Context(), userID, *payload.ParentTodoID)
		if err != nil {
			logger.Error().Err(err).Msg("parent todo validation failed ")
			return nil, err
		}

		if !parentTodo.CanHaveChildren() {
			err := errs.NewBadRequestError("parent todo cannot have children (subtasks can't have subtasks)", false, nil, nil, nil)
			logger.Warn().Msg("parent todo cannot have children")
			return nil, err
		}
	}

	if payload.CategoryID != nil {
		_, err := s.categoryRepo.GetCategoryByID(ctx.Request().Context(), userID, *payload.CategoryID)
		if err != nil {
			logger.Error().Err(err).Msg("category validation failed")
			return nil, err

		}
	}

	todoItem, err := s.todoRepo.CreateTodo(context.Background(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create todo")
		return nil, err

	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "tod_created").
		Str("todo_id", todoItem.ID.String()).
		Str("title", todoItem.Title).
		Str("category_id", func() string {
			if todoItem.CategoryID != nil {
				return todoItem.CategoryID.String()
			}

			return ""
		}()).
		Str("priority", string(todoItem.Priority)).
		Msg("Todo creted successfullyt")

	return todoItem, nil
}
