package handler

import (
	"net/http"

	"github.com/C0deNe0/go-tasker/internal/middleware"
	"github.com/C0deNe0/go-tasker/internal/model/comment"
	"github.com/C0deNe0/go-tasker/internal/server"
	"github.com/C0deNe0/go-tasker/internal/service"
	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	Handler
	commentService *service.CommentService
}

func NewCommentHandler(s *server.Server, commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		Handler:        NewHandler(s),
		commentService: commentService,
	}
}

func (h *CommentHandler) AddComment(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *comment.AddCommentPayload) (*comment.Comment, error) {
			userID := middleware.GetUserID(c)
			return h.commentService.AddComment(c, userID, payload.TodoID, payload)
		},
		http.StatusCreated,
		&comment.AddCommentPayload{},
	)(c)
}

func (h *CommentHandler) GetCommentsByTodoID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *comment.GetCommentByTodoIDPayload) ([]comment.Comment, error) {
			userID := middleware.GetUserID(c)
			return h.commentService.GetCommentsByTodoID(c, userID, payload.TodoID)
		},
		http.StatusOK,
		&comment.GetCommentByTodoIDPayload{},
	)(c)
}

func (h *CommentHandler) UpdateComment(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *comment.UpdateCommentPayload) (*comment.Comment, error) {
			userID := middleware.GetUserID(c)
			return h.commentService.UpdateComment(c, userID, payload.ID, payload.Content)
		},
		http.StatusOK,
		&comment.UpdateCommentPayload{},
	)(c)
}

func (h *CommentHandler) DeleteComment(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *comment.DeleteCommentPayload) error {
			userID := middleware.GetUserID(c)
			return h.commentService.DeleteComment(c, userID, payload.ID)
		},
		http.StatusNoContent,
		&comment.DeleteCommentPayload{},
	)(c)
}
