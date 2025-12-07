package todo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestGetTodos(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	c := e.NewContext(req, w)

	handler := NewTodoHandler()
	err := handler.getTodos(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "Get todos", w.Body.String())
}

func TestGetTodo(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
	w := httptest.NewRecorder()

	c := e.NewContext(req, w)

	handler := NewTodoHandler()
	err := handler.getTodo(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "Get todo", w.Body.String())
}
