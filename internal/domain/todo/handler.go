package todo

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct{}

func (h *Handler) RegisterTodoRoutes(e *echo.Echo) {
	g := e.Group("/todos")
	g.GET("", h.getTodos)
	g.GET(":id", h.getTodo)
}

func (h *Handler) getTodos(c echo.Context) error {
	return c.String(http.StatusOK, "Get todos")
}

func (h *Handler) getTodo(c echo.Context) error {
	return c.String(http.StatusOK, "Get todo")
}

func NewTodoHandler() *Handler {
	return &Handler{}
}
