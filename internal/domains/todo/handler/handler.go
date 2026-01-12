package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct{}

func (h *Handler) RegisterTodoRoutes(e *echo.Echo) {
	g := e.Group("/todos")
	g.POST("", h.createTodo)
	g.GET("", h.getTodos)
	g.GET(":id", h.getTodo)
	g.PUT(":id", h.updateTodo)
	g.DELETE(":id", h.deleteTodo)
}

func NewTodoHandler() *Handler {
	return &Handler{}
}

func (h *Handler) createTodo(c echo.Context) error { return c.NoContent(http.StatusCreated) }

func (h *Handler) getTodos(c echo.Context) error { return c.String(http.StatusOK, "Get todos") }

func (h *Handler) getTodo(c echo.Context) error { return c.String(http.StatusOK, "Get todo") }

func (h *Handler) updateTodo(c echo.Context) error { return c.NoContent(http.StatusNoContent) }

func (h *Handler) deleteTodo(c echo.Context) error { return c.NoContent(http.StatusNoContent) }
