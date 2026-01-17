package handler

import (
	"context"
	"echo-sample2/internal/todo/dto"
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

type MockTodoService struct {
	mock.Mock
}

func (m *MockTodoService) Create(ctx context.Context, input dto.CreateTodoInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func TestTodoHandler(t *testing.T) {
	setup := func() (*echo.Echo, *TodoHandler, *MockTodoService) {
		e := echo.New()
		e.Validator = &customValidator.CustomValidator{Validator: validator.New()}
		mockService := new(MockTodoService)
		h := NewTodoHandler(mockService)
		h.RegisterTodoRoutes(e)
		return e, h, mockService
	}

	t.Run("CreateTodo", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _, mockService := setup()
			input := dto.CreateTodoInput{
				Title:       "test",
				Description: "description",
			}
			mockService.On("Create", mock.Anything, input).Return(nil)

			req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(`{"title": "test", "description": "description"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusCreated, rec.Code)
			require.Equal(t, "", rec.Body.String())
			mockService.AssertExpectations(t)
		})

		t.Run("バリデーションエラー", func(t *testing.T) {
			e, _, mockService := setup()
			req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(`{}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
			mockService.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
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
