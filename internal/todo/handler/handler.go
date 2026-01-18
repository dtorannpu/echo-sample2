package handler

import (
	"context"
	"echo-sample2/internal/todo/application/usecase/create"
	"echo-sample2/internal/todo/application/usecase/list"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateTodoUseCase interface {
	Execute(ctx context.Context, command create.Command) error
}

type GetTodosUseCase interface {
	Execute(ctx context.Context) (*list.Result, error)
}

type TodoHandler struct {
	createTodoUseCase CreateTodoUseCase
	getTodosUseCase   GetTodosUseCase
}

type createTodoRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type todoResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type getTodosResponse struct {
	Todos []*todoResponse `json:"todos"`
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

func NewTodoHandler(createTodoUseCase CreateTodoUseCase, getTodosUseCase GetTodosUseCase) *TodoHandler {
	return &TodoHandler{createTodoUseCase: createTodoUseCase, getTodosUseCase: getTodosUseCase}
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

	command := create.Command{
		Title:       req.Title,
		Description: req.Description,
	}

	err := h.createTodoUseCase.Execute(c.Request().Context(), command)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (h *TodoHandler) getTodos(c echo.Context) error {
	todoRes, err := h.getTodosUseCase.Execute(c.Request().Context())
	if err != nil {
		return err
	}

	res := make([]*todoResponse, len(todoRes.Todos))

	for i, todo := range todoRes.Todos {
		res[i] = &todoResponse{
			ID:          todo.ID.String(),
			Title:       todo.Title,
			Description: todo.Description,
		}
	}

	return c.JSON(http.StatusOK, &getTodosResponse{Todos: res})
}

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
