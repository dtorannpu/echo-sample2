package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TodoHandler struct{}

type createTodoRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type getTodoRequest struct {
	ID int64 `param:"id" validate:"required"`
}

type updateTodoRequest struct {
	ID          int64  `param:"id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type deleteTodoRequest struct {
	ID int64 `param:"id" validate:"required"`
}

func NewTodoHandler() *TodoHandler {
	return &TodoHandler{}
}

func (h *TodoHandler) RegisterTodoRoutes(e *echo.Echo) {
	g := e.Group("/todos")
	g.POST("", h.createTodo)
	g.GET("", h.getTodos)
	g.GET("/:id", h.getTodo)
	g.PUT("/:id", h.updateTodo)
	g.DELETE("/:id", h.deleteTodo)
}

func (h *TodoHandler) createTodo(c echo.Context) error {
	req := new(createTodoRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (h *TodoHandler) getTodos(c echo.Context) error { return c.String(http.StatusOK, "Get todos") }

func (h *TodoHandler) getTodo(c echo.Context) error {
	req := new(getTodoRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "Get todo")
}

func (h *TodoHandler) updateTodo(c echo.Context) error {
	req := new(updateTodoRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *TodoHandler) deleteTodo(c echo.Context) error {
	req := new(deleteTodoRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
