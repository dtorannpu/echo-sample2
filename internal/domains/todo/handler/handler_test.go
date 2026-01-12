package handler

import (
	customValidator "echo-sample2/internal/validator"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestTodoHandler(t *testing.T) {
	setup := func() (*echo.Echo, *TodoHandler) {
		e := echo.New()
		e.Validator = &customValidator.CustomValidator{Validator: validator.New()}
		h := NewTodoHandler()
		h.RegisterTodoRoutes(e)
		return e, h
	}

	t.Run("CreateTodo", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _ := setup()
			req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(`{"title": "test", "description": "description"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusCreated, rec.Code)
			require.Equal(t, "", rec.Body.String())
		})

		t.Run("バリデーションエラー", func(t *testing.T) {
			e, _ := setup()
			req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(`{}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})

	t.Run("GetTodos", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _ := setup()
			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
			require.Equal(t, "Get todos", rec.Body.String())
		})
	})

	t.Run("GetTodo", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _ := setup()
			req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
			require.Equal(t, "Get todo", rec.Body.String())
		})

		t.Run("IDが数字じゃない場合", func(t *testing.T) {
			e, _ := setup()
			req := httptest.NewRequest(http.MethodGet, "/todos/a", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues("a")
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})

	t.Run("UpdateTodo", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _ := setup()
			req := httptest.NewRequest(http.MethodPut, "/todos/1", strings.NewReader(`{"title": "test", "description": "description"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusNoContent, rec.Code)
			require.Equal(t, "", rec.Body.String())
		})

		t.Run("IDが数字じゃない場合", func(t *testing.T) {
			e, _ := setup()
			req := httptest.NewRequest(http.MethodPut, "/todos/a", strings.NewReader(`{"title": "test", "description": "description"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues("a")
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
		})

		t.Run("バリデーションエラー", func(t *testing.T) {
			e, _ := setup()
			req := httptest.NewRequest(http.MethodPut, "/todos/a", strings.NewReader(`{}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues("a")
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})

	t.Run("GetTodos", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _ := setup()
			req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusNoContent, rec.Code)
			require.Equal(t, "", rec.Body.String())
		})

		t.Run("IDが数字じゃない場合", func(t *testing.T) {
			e, _ := setup()
			req := httptest.NewRequest(http.MethodDelete, "/todos/a", strings.NewReader(`{}`))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues("a")
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})
}
