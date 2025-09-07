package v1

import (
	"github.com/C0deNe0/go-tasker/internal/handler"
	"github.com/C0deNe0/go-tasker/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerCategoryRoutes(r *echo.Group, h *handler.CategoryHandler, auth *middleware.AuthMiddleware) {

	categories := r.Group("/categories")
	categories.Use(auth.RequireAuth)

	categories.POST("", h.CreateCategory)
	categories.GET("", h.GetCategories)

	dynamicCategory := categories.Group("/:id")
	dynamicCategory.PATCH("", h.UpdateCategory)
	dynamicCategory.DELETE("", h.DeleteCategory)
}
