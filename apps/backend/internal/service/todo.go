package service

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/C0deNe0/go-tasker/internal/errs"
	"github.com/C0deNe0/go-tasker/internal/lib/aws"
	"github.com/C0deNe0/go-tasker/internal/middleware"
	"github.com/C0deNe0/go-tasker/internal/model"
	"github.com/C0deNe0/go-tasker/internal/model/todo"
	"github.com/C0deNe0/go-tasker/internal/repository"
	"github.com/C0deNe0/go-tasker/internal/server"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type TodoService struct {
	server       *server.Server
	todoRepo     *repository.TodoRepository
	categoryRepo *repository.CategoryRepository
	awsClient    *aws.AWS
}

func NewTodoService(server *server.Server, todoRepo *repository.TodoRepository, categroyRepo *repository.CategoryRepository, awsClient *aws.AWS) *TodoService {
	return &TodoService{
		server:       server,
		todoRepo:     todoRepo,
		categoryRepo: categroyRepo,
		awsClient:    awsClient,
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

func (s *TodoService) GetTodoByID(ctx echo.Context, userID string, todoID uuid.UUID) (*todo.PopulatedTodo, error) {
	logger := middleware.GetLogger(ctx)

	todoItem, err := s.todoRepo.GetTodoByID(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch by ID")
		return nil, err

	}
	return todoItem, nil
}

func (s *TodoService) GetTodos(ctx echo.Context, userID string, query *todo.GetTodosQuery) (*model.PaginatedResponse[todo.PopulatedTodo], error) {
	logger := middleware.GetLogger(ctx)

	result, err := s.todoRepo.GetTodos(ctx.Request().Context(), userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch todos")
		return nil, err

	}

	return result, nil
}

func (s *TodoService) UpdateTodo(ctx echo.Context, userID string, payload *todo.UpdateTodoPayload) (*todo.Todo, error) {
	logger := middleware.GetLogger(ctx)

	// Validate parent todo exists and belongs to user (if provided)
	if payload.ParentTodoID != nil {
		parentTodo, err := s.todoRepo.CheckTodoExists(ctx.Request().Context(), userID, *payload.ParentTodoID)
		if err != nil {
			logger.Error().Err(err).Msg("parent todo validation failed")
			return nil, err
		}

		if parentTodo.ID == payload.ID {
			err := errs.NewBadRequestError("Todo cannot be its own parent", false, nil, nil, nil)
			logger.Warn().Msg("todo cannot be its own parent")
			return nil, err
		}

		if !parentTodo.CanHaveChildren() {
			err := errs.NewBadRequestError("Parent todo cannot have children (subtasks can't have subtasks)", false, nil, nil, nil)
			logger.Warn().Msg("parent todo cannot have children")
			return nil, err
		}

		logger.Debug().Msg("parent todo validation passed")
	}

	// Validate category exists and belongs to user (if provided)
	if payload.CategoryID != nil {
		_, err := s.categoryRepo.GetCategoryByID(ctx.Request().Context(), userID, *payload.CategoryID)
		if err != nil {
			logger.Error().Err(err).Msg("category validation failed")
			return nil, err
		}

		logger.Debug().Msg("category validation passed")
	}

	updatedTodo, err := s.todoRepo.UpdateTodo(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update todo")
		return nil, err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "todo_updated").
		Str("todo_id", updatedTodo.ID.String()).
		Str("title", updatedTodo.Title).
		Str("category_id", func() string {
			if updatedTodo.CategoryID != nil {
				return updatedTodo.CategoryID.String()
			}
			return ""
		}()).
		Str("priority", string(updatedTodo.Priority)).
		Str("status", string(updatedTodo.Status)).
		Msg("Todo updated successfully")

	return updatedTodo, nil
}

func (s *TodoService) DeleteTodo(ctx echo.Context, userID string, todoID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	err := s.todoRepo.DeleteTodo(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete todo")
		return err

	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event ", "todo_deleted").
		Str("todo_id", todoID.String()).
		Msg("todo deleted successfully")

	return nil

}

func (s *TodoService) GetTodoStats(ctx echo.Context, userID string) (*todo.TodoStats, error) {
	logger := middleware.GetLogger(ctx)

	stats, err := s.todoRepo.GetTodoStats(ctx.Request().Context(), userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch stats")
		return nil, err
	}

	return stats, nil

}

func (s *TodoService) UploadTodoAttachment(ctx echo.Context, userID string, todoID string, file *multipart.FileHeader) (*todo.TodoAttachment, error) {
	logger := middleware.GetLogger(ctx)

	//parse todo UUID
	todoUUID, err := uuid.Parse(todoID)
	if err != nil {
		logger.Error().Err(err).Msg("invalid todo ID")
		return nil, errs.NewBadRequestError("invalid todo id", false, nil, nil, nil)
	}

	//verify exist or not
	_, err = s.todoRepo.CheckTodoExists(ctx.Request().Context(), userID, todoUUID)
	if err != nil {
		logger.Error().Err(err).Msg("todo validation failed")
		return nil, err
	}

	//open upload file
	src, err := file.Open()
	if err != nil {
		logger.Error().Err(err).Msg("failed to open the file")
		return nil, errs.NewBadRequestError("failed to open uploaded file", false, nil, nil, nil)
	}

	defer src.Close()

	//uploading to S3

	s3Key, err := s.awsClient.S3.UploadFile(
		ctx.Request().Context(),
		s.server.Config.AWS.UploadBucket,
		"todos/attachments/"+file.Filename,
		src,
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to upload file to s3")
		return nil, errors.Wrap(err, " failed to upload file")
	}

	//Delete Mime type
	src, err = file.Open()
	if err != nil {
		logger.Error().Err(err).Msg("failed to reopen for MIME detection")
		return nil, errs.NewBadRequestError("failed to process file", false, nil, nil, nil)
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		logger.Error().Err(err).Msg("failed to read for MIME detection")
		return nil, errs.NewBadRequestError("failed to process file", false, nil, nil, nil)

	}

	mimeType := http.DetectContentType(buffer)

	attachment, err := s.todoRepo.UploadTodoAttachment(
		ctx.Request().Context(),
		todoUUID,
		userID,
		s3Key,
		file.Filename,
		file.Size,
		mimeType,
	)

	if err != nil {
		logger.Error().Err(err).Msg("failed to create attachment record")
		return nil, err
	}

	logger.Info().
		Str("attachmentID", attachment.ID.String()).
		Str("s3_key", s3Key).
		Msg("uploaded todo attachment")

	return attachment, nil
}

func (s *TodoService) DeleteTodoAttachment(ctx echo.Context, userID string, todoID uuid.UUID, attachmentID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	_, err := s.todoRepo.GetTodoAttachment(ctx.Request().Context(), todoID, attachmentID)
	if err != nil {
		logger.Error().Err(err).Msg("todo validation failed")
		return err
	}

	//get attachment details for s3 deletion
	attachment, err := s.todoRepo.GetTodoAttachment(
		ctx.Request().Context(),
		todoID,
		attachmentID,
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get attachment detailes")
		return err
	}

	//delete attachment record
	err = s.todoRepo.DeleteTodoAttachment(ctx.Request().Context(),
		todoID, attachmentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete attachment record")
		return err
	}

	//delete the object without stoping for the response
	go func() {
		err := s.awsClient.S3.DeleteObject(ctx.Request().Context(), s.server.Config.AWS.UploadBucket, attachment.DownloadKey)

		if err != nil {
			s.server.Logger.Error().Err(err).Str("s3_key", attachment.DownloadKey).Msg("failed to delete attachment from s3")
		}

	}()
	logger.Info().Msg("deleted todo message")
	return nil
}

func (s *TodoService) GetAttachmentPresignedURL(ctx echo.Context, userID string, todoID uuid.UUID, attachmentID uuid.UUID) (string, error) {
	logger := middleware.GetLogger(ctx)

	_, err := s.todoRepo.CheckTodoExists(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("todo validation failed")
		return "", err
	}

	attachment, err := s.todoRepo.GetTodoAttachment(
		ctx.Request().Context(),
		todoID,
		attachmentID,
	)

	if err != nil {

		logger.Error().Err(err).Msg("failed to get attachment file")
		return "", err
	}

	//generate the PResigned URL
	url, err := s.awsClient.S3.CreatePresignedUrl(
		ctx.Request().Context(),
		s.server.Config.AWS.UploadBucket,
		attachment.DownloadKey,
	)

	if err != nil {
		logger.Error().Err(err).Msg("failed to generate presigned URL")
		return "", err
	}

	return url, nil
}
