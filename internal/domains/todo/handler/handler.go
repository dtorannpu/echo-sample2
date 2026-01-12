package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TodoHandler struct{}

type getTodoParams struct {
	ID int64 `param:"id"`
}

type updateTodoParams struct {
	ID int64 `param:"id"`
}

type deleteTodoParams struct {
	ID int64 `param:"id"`
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

func (h *TodoHandler) createTodo(c echo.Context) error { return c.NoContent(http.StatusCreated) }

func (h *TodoHandler) getTodos(c echo.Context) error { return c.String(http.StatusOK, "Get todos") }

func (h *TodoHandler) getTodo(c echo.Context) error {
	params := new(getTodoParams)

	if err := c.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "Get todo")
}

func (h *TodoHandler) updateTodo(c echo.Context) error {
	params := new(updateTodoParams)

	if err := c.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *TodoHandler) deleteTodo(c echo.Context) error {
	params := new(deleteTodoParams)

	if err := c.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
