package handler

import (
	"context"
	"echo-sample2/internal/todo/application/usecase"
	"echo-sample2/internal/todo/domain"
	customValidator "echo-sample2/internal/validator"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
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

type MockGetTodosUseCase struct {
	mock.Mock
}

func (m *MockGetTodosUseCase) Execute(ctx context.Context) (*usecase.GetTodosResult, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.GetTodosResult), args.Error(1)
}

func TestTodoHandler(t *testing.T) {
	setup := func() (*echo.Echo, *TodoHandler, *MockCreateTodoUseCase, *MockGetTodosUseCase) {
		e := echo.New()
		e.Validator = &customValidator.CustomValidator{Validator: validator.New()}
		mockCreateUseCase := new(MockCreateTodoUseCase)
		mockGetUseCase := new(MockGetTodosUseCase)
		h := NewTodoHandler(mockCreateUseCase, mockGetUseCase)
		h.RegisterTodoRoutes(e)
		return e, h, mockCreateUseCase, mockGetUseCase
	}

	t.Run("CreateTodo", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _, mockUseCase, _ := setup()
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
			e, _, mockUseCase, _ := setup()
			req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(`{}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
			mockUseCase.AssertNotCalled(t, "Execute", mock.Anything, mock.Anything)
		})

		t.Run("ユースケースでエラーが発生したパターン", func(t *testing.T) {
			e, _, mockUseCase, _ := setup()
			command := usecase.CreateTodoCommand{
				Title:       "test",
				Description: "description",
			}
			mockUseCase.On("Execute", mock.Anything, command).Return(fmt.Errorf("usecase error"))

			req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(`{"title": "test", "description": "description"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusInternalServerError, rec.Code)
			mockUseCase.AssertExpectations(t)
		})
	})

	t.Run("GetTodos", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _, _, mockGetUseCase := setup()
			mockGetUseCase.On("Execute", mock.Anything).Return(&usecase.GetTodosResult{Todos: []*usecase.TodoResult{}}, nil)
			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
			require.JSONEq(t, `{"todos":[]}`, rec.Body.String())
		})

		t.Run("データが返却されるパターン", func(t *testing.T) {
			e, _, _, mockGetUseCase := setup()
			id := uuid.New()
			mockGetUseCase.On("Execute", mock.Anything).Return(&usecase.GetTodosResult{
				Todos: []*usecase.TodoResult{
					{
						ID:          domain.TodoID(id),
						Title:       "test title",
						Description: "test description",
					},
				},
			}, nil)
			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
			expectedBody := fmt.Sprintf(`{"todos":[{"id":"%s","title":"test title","description":"test description"}]}`, id.String())
			require.JSONEq(t, expectedBody, rec.Body.String())
		})

		t.Run("ユースケースでエラーが発生したパターン", func(t *testing.T) {
			e, _, _, mockGetUseCase := setup()
			mockGetUseCase.On("Execute", mock.Anything).Return(nil, fmt.Errorf("usecase error"))
			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusInternalServerError, rec.Code)
		})
	})

	t.Run("GetTodo", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _, _, _ := setup()
			req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
			require.Equal(t, "Get todo", rec.Body.String())
		})

		t.Run("IDが数字じゃない場合", func(t *testing.T) {
			e, _, _, _ := setup()
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
			e, _, _, _ := setup()
			req := httptest.NewRequest(http.MethodPut, "/todos/1", strings.NewReader(`{"title": "test", "description": "description"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusNoContent, rec.Code)
			require.Equal(t, "", rec.Body.String())
		})

		t.Run("IDが数字じゃない場合", func(t *testing.T) {
			e, _, _, _ := setup()
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
			e, _, _, _ := setup()
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
			e, _, _, _ := setup()
			req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusNoContent, rec.Code)
			require.Equal(t, "", rec.Body.String())
		})

		t.Run("IDが数字じゃない場合", func(t *testing.T) {
			e, _, _, _ := setup()
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
