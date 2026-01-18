package handler

import (
	"context"
	"echo-sample2/internal/todo/application/usecase/create"
	"echo-sample2/internal/todo/application/usecase/get"
	"echo-sample2/internal/todo/application/usecase/list"
	"echo-sample2/internal/todo/domain"
	domainError "echo-sample2/internal/todo/domain/errors"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateUseCase interface {
	Execute(ctx context.Context, command create.Command) error
}

type ListUseCase interface {
	Execute(ctx context.Context) (*list.Result, error)
}

type GetUseCase interface {
	Execute(ctx context.Context, command get.Command) (*get.Todo, error)
}

type TodoHandler struct {
	createUseCase CreateUseCase
	listUseCase   ListUseCase
	getUseCase    GetUseCase
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
	ID string `param:"id" validate:"required"`
}

type updateTodoRequest struct {
	ID          string `param:"id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type deleteTodoRequest struct {
	ID string `param:"id" validate:"required"`
}

func NewTodoHandler(createTodoUseCase CreateUseCase, listUseCase ListUseCase, getUseCase GetUseCase) *TodoHandler {
	return &TodoHandler{
		createUseCase: createTodoUseCase,
		listUseCase:   listUseCase,
		getUseCase:    getUseCase,
	}
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

	err := h.createUseCase.Execute(c.Request().Context(), command)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (h *TodoHandler) getTodos(c echo.Context) error {
	todoRes, err := h.listUseCase.Execute(c.Request().Context())
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

	id, err := domain.NewTodoIDFromString(req.ID)
	var invalidErr *domainError.ErrInvalidTodoID
	if errors.As(err, &invalidErr) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid todo ID")
	}
	if err != nil {
		return err
	}

	res, err := h.getUseCase.Execute(c.Request().Context(), get.Command{ID: *id})
	if err != nil {
		return err
	}
	if res == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Todo not found")
	}
	return c.JSON(http.StatusOK, &todoResponse{
		ID:          res.ID.String(),
		Title:       res.Title,
		Description: res.Description,
	})
}

func (h *TodoHandler) updateTodo(c echo.Context) error {
	req := new(updateTodoRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	_, err := domain.NewTodoIDFromString(req.ID)
	var invalidErr *domainError.ErrInvalidTodoID
	if errors.As(err, &invalidErr) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid todo ID")
	}
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *TodoHandler) deleteTodo(c echo.Context) error {
	req := new(deleteTodoRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := domain.NewTodoIDFromString(req.ID)
	var invalidErr *domainError.ErrInvalidTodoID
	if errors.As(err, &invalidErr) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid todo ID")
	}
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
