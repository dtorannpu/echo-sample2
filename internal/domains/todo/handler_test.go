package todo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestCreateTodo(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewTodoHandler()

	require.NoError(t, h.createTodo(c))
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, "", rec.Body.String())
}

func TestGetTodos(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewTodoHandler()

	require.NoError(t, h.getTodos(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "Get todos", rec.Body.String())
}

func TestGetTodo(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewTodoHandler()

	require.NoError(t, h.getTodo(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "Get todo", rec.Body.String())
}

func TestUpdateTodo(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/todos/1", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewTodoHandler()

	require.NoError(t, h.updateTodo(c))
	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Equal(t, "", rec.Body.String())
}

func TestDeleteTodo(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewTodoHandler()

	require.NoError(t, h.updateTodo(c))
	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Equal(t, "", rec.Body.String())
}
