package v1

import (
	"github.com/C0deNe0/go-tasker/internal/handler"
	"github.com/C0deNe0/go-tasker/internal/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterV1Routes(routes *echo.Group, handlers *handler.Handlers, middleware *middleware.Middlewares) {
	//register todo route
	registerTodoRoutes(routes, handlers.Todo, handlers.Comment, middleware.Auth)
	//category
	registerCategoryRoutes(routes, handlers.Category, middleware.Auth)
	//comments
	registerCommentRoutes(routes, handlers.Comment, middleware.Auth)
}
