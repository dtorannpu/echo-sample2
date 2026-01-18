package handler

import (
	"context"
	"echo-sample2/internal/todo/application/usecase"
	customValidator "echo-sample2/internal/validator"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCreateTodoUseCase struct {
	mock.Mock
}

func (m *MockCreateTodoUseCase) Execute(ctx context.Context, command usecase.CreateTodoCommand) error {
	args := m.Called(ctx, command)
	return args.Error(0)
}

func TestTodoHandler(t *testing.T) {
	setup := func() (*echo.Echo, *TodoHandler, *MockCreateTodoUseCase) {
		e := echo.New()
		e.Validator = &customValidator.CustomValidator{Validator: validator.New()}
		mockUseCase := new(MockCreateTodoUseCase)
		h := NewTodoHandler(mockUseCase)
		h.RegisterTodoRoutes(e)
		return e, h, mockUseCase
	}

	t.Run("CreateTodo", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _, mockUseCase := setup()
			command := usecase.CreateTodoCommand{
				Title:       "test",
				Description: "description",
			}
			mockUseCase.On("Execute", mock.Anything, command).Return(nil)

			req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(`{"title": "test", "description": "description"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusCreated, rec.Code)
			require.Equal(t, "", rec.Body.String())
			mockUseCase.AssertExpectations(t)
		})

		t.Run("バリデーションエラー", func(t *testing.T) {
			e, _, mockUseCase := setup()
			req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(`{}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
			mockUseCase.AssertNotCalled(t, "Execute", mock.Anything, mock.Anything)
		})
	})

	t.Run("GetTodos", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _, _ := setup()
			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
			require.Equal(t, "Get todos", rec.Body.String())
		})
	})

	t.Run("GetTodo", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _, _ := setup()
			req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
			require.Equal(t, "Get todo", rec.Body.String())
		})

		t.Run("IDが数字じゃない場合", func(t *testing.T) {
			e, _, _ := setup()
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
			e, _, _ := setup()
			req := httptest.NewRequest(http.MethodPut, "/todos/1", strings.NewReader(`{"title": "test", "description": "description"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusNoContent, rec.Code)
			require.Equal(t, "", rec.Body.String())
		})

		t.Run("IDが数字じゃない場合", func(t *testing.T) {
			e, _, _ := setup()
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
			e, _, _ := setup()
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

	t.Run("DeleteTodo", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _, _ := setup()
			req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusNoContent, rec.Code)
			require.Equal(t, "", rec.Body.String())
		})

		t.Run("IDが数字じゃない場合", func(t *testing.T) {
			e, _, _ := setup()
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
