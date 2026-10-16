package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/poojanbhardwaj/task-manager/apps/backend/internal/handler"
	"github.com/poojanbhardwaj/task-manager/apps/backend/internal/middleware"
)

func registerTodoRoutes(r *echo.Group, h *handler.TodoHandler, ch *handler.CommentHandler, auth *middleware.AuthMiddleware) {
	// Todo operations
	todos := r.Group("/todos")
	todos.Use(auth.RequireAuth)

	// Collection operations
	todos.POST("", h.CreateTodo)
	todos.GET("", h.GetTodos)
	todos.GET("/stats", h.GetTodoStats)

	// Individual todo operations
	dynamicTodo := todos.Group("/:id")
	dynamicTodo.GET("", h.GetTodoByID)
	dynamicTodo.PATCH("", h.UpdateTodo)
	dynamicTodo.DELETE("", h.DeleteTodo)
	dynamicTodo.POST("/attachments", h.UploadTodoAttachment)
	dynamicTodo.DELETE("/attachments/:attachmentId", h.DeleteTodoAttachment)
	dynamicTodo.GET("/attachments/:attachmentId/download", h.GetAttachmentPresignedURL)

	// Todo comments
	todoComments := dynamicTodo.Group("/comments")
	todoComments.POST("", ch.AddComment)
	todoComments.GET("", ch.GetCommentsByTodoID)
}
