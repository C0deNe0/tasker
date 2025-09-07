package v1

import (
	"github.com/C0deNe0/go-tasker/internal/handler"
	"github.com/C0deNe0/go-tasker/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerTodoRoutes(r *echo.Group, h *handler.TodoHandler, ch *handler.CommentHandler, auth *middleware.AuthMiddleware) {

	//todo opertn
	todos := r.Group("/todos")
	todos.Use(auth.RequireAuth)

	//collection operations
	todos.POST("", h.CreateTodo)
	todos.GET("", h.GetTodos)
	todos.GET("/stats", h.GetTodoStats)

	dynamicTodo := todos.Group("/:id")
	dynamicTodo.GET("", h.GetTodoByID)
	dynamicTodo.PATCH("", h.UpdateTodo)
	dynamicTodo.DELETE("", h.DeleteTodo)

	//commetns
	todoComments := dynamicTodo.Group("/comments")
	todoComments.PUT("", ch.AddComment)
	todoComments.GET("", ch.GetCommentsByTodoID)
}
