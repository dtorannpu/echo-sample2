package handler

import (
	"context"
	"echo-sample2/internal/todo/application/usecase/create"
	"echo-sample2/internal/todo/application/usecase/delete"
	"echo-sample2/internal/todo/application/usecase/get"
	"echo-sample2/internal/todo/application/usecase/list"
	"echo-sample2/internal/todo/application/usecase/update"
	"echo-sample2/internal/todo/domain"
	domainError "echo-sample2/internal/todo/domain/errors"
	customValidator "echo-sample2/internal/validator"
	"fmt"
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

func (m *MockCreateTodoUseCase) Execute(ctx context.Context, command create.Command) error {
	args := m.Called(ctx, command)
	return args.Error(0)
}

type MockGetTodosUseCase struct {
	mock.Mock
}

func (m *MockGetTodosUseCase) Execute(ctx context.Context) (*list.Result, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*list.Result), args.Error(1)
}

type MockGetTodoUseCase struct {
	mock.Mock
}

func (m *MockGetTodoUseCase) Execute(ctx context.Context, command get.Command) (*get.Todo, error) {
	args := m.Called(ctx, command)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*get.Todo), args.Error(1)
}

type MockUpdateTodoUseCase struct {
	mock.Mock
}

func (m *MockUpdateTodoUseCase) Execute(ctx context.Context, command update.Command) error {
	args := m.Called(ctx, command)
	return args.Error(0)
}

type MockDeleteTodoUseCase struct {
	mock.Mock
}

func (m *MockDeleteTodoUseCase) Execute(ctx context.Context, command delete.Command) error {
	args := m.Called(ctx, command)
	return args.Error(0)
}

func TestTodoHandler(t *testing.T) {
	setup := func() (*echo.Echo, *TodoHandler, *MockCreateTodoUseCase, *MockGetTodosUseCase, *MockGetTodoUseCase, *MockUpdateTodoUseCase, *MockDeleteTodoUseCase) {
		e := echo.New()
		e.Validator = &customValidator.CustomValidator{Validator: validator.New()}
		mockCreateUseCase := new(MockCreateTodoUseCase)
		mockGetTodosUseCase := new(MockGetTodosUseCase)
		mockGetTodoUseCase := new(MockGetTodoUseCase)
		mockUpdateUseCase := new(MockUpdateTodoUseCase)
		mockDeleteUseCase := new(MockDeleteTodoUseCase)
		h := NewTodoHandler(mockCreateUseCase, mockGetTodosUseCase, mockGetTodoUseCase, mockUpdateUseCase, mockDeleteUseCase)
		h.RegisterTodoRoutes(e)
		return e, h, mockCreateUseCase, mockGetTodosUseCase, mockGetTodoUseCase, mockUpdateUseCase, mockDeleteUseCase
	}

	t.Run("CreateTodo", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _, mockUseCase, _, _, _, _ := setup()
			command := create.Command{
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
			e, _, _, _, _, _, _ := setup()
			req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(`{}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
		})

		t.Run("ユースケースでエラーが発生したパターン", func(t *testing.T) {
			e, _, mockUseCase, _, _, _, _ := setup()
			command := create.Command{
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
			e, _, _, mockGetUseCase, _, _, _ := setup()
			mockGetUseCase.On("Execute", mock.Anything).Return(&list.Result{Todos: []*list.Todo{}}, nil)
			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
			require.JSONEq(t, `{"todos":[]}`, rec.Body.String())
		})

		t.Run("データが返却されるパターン", func(t *testing.T) {
			e, _, _, mockGetUseCase, _, _, _ := setup()
			id, err := domain.NewTodoID()
			require.NoError(t, err)
			mockGetUseCase.On("Execute", mock.Anything).Return(&list.Result{
				Todos: []*list.Todo{
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
			e, _, _, mockGetUseCase, _, _, _ := setup()
			mockGetUseCase.On("Execute", mock.Anything).Return(nil, fmt.Errorf("usecase error"))
			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusInternalServerError, rec.Code)
		})
	})

	t.Run("GetTodo", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _, _, _, mockGetTodoUseCase, _, _ := setup()
			id, err := domain.NewTodoID()
			require.NoError(t, err)
			todoID, _ := domain.NewTodoIDFromString(id.String())
			mockGetTodoUseCase.On("Execute", mock.Anything, get.Command{ID: *todoID}).Return(&get.Todo{
				ID:          *todoID,
				Title:       "test title",
				Description: "test description",
			}, nil)

			req := httptest.NewRequest(http.MethodGet, "/todos/"+id.String(), nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
			expectedBody := fmt.Sprintf(`{"id":"%s","title":"test title","description":"test description"}`, id.String())
			require.JSONEq(t, expectedBody, rec.Body.String())
			mockGetTodoUseCase.AssertExpectations(t)
		})

		t.Run("IDがUUIDじゃない場合", func(t *testing.T) {
			e, _, _, _, _, _, _ := setup()
			req := httptest.NewRequest(http.MethodGet, "/todos/invalid-uuid", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
			require.JSONEq(t, `{"message":"Invalid todo ID"}`, rec.Body.String())
		})

		t.Run("データが存在しない場合", func(t *testing.T) {
			e, _, _, _, mockGetTodoUseCase, _, _ := setup()
			id, err := domain.NewTodoID()
			require.NoError(t, err)
			todoID, _ := domain.NewTodoIDFromString(id.String())
			mockGetTodoUseCase.On("Execute", mock.Anything, get.Command{ID: *todoID}).Return(nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/todos/"+id.String(), nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusNotFound, rec.Code)
			require.JSONEq(t, `{"message":"Todo not found"}`, rec.Body.String())
			mockGetTodoUseCase.AssertExpectations(t)
		})
	})

	t.Run("UpdateTodo", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _, _, _, _, mockUpdateUseCase, _ := setup()
			id, err := domain.NewTodoID()
			require.NoError(t, err)
			todoID, _ := domain.NewTodoIDFromString(id.String())
			command := update.Command{
				ID:          *todoID,
				Title:       "test",
				Description: "description",
			}
			mockUpdateUseCase.On("Execute", mock.Anything, command).Return(nil)

			req := httptest.NewRequest(http.MethodPut, "/todos/"+id.String(), strings.NewReader(`{"title": "test", "description": "description"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusNoContent, rec.Code)
			require.Equal(t, "", rec.Body.String())
			mockUpdateUseCase.AssertExpectations(t)
		})

		t.Run("IDがUUIDじゃない場合", func(t *testing.T) {
			e, _, _, _, _, _, _ := setup()
			req := httptest.NewRequest(http.MethodPut, "/todos/invalid-uuid", strings.NewReader(`{"title": "test", "description": "description"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
			require.JSONEq(t, `{"message":"Invalid todo ID"}`, rec.Body.String())
		})

		t.Run("バリデーションエラー", func(t *testing.T) {
			e, _, _, _, _, _, _ := setup()
			id, err := domain.NewTodoID()
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPut, "/todos/"+id.String(), strings.NewReader(`{}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
		})

		t.Run("ユースケースでエラーが発生したパターン", func(t *testing.T) {
			e, _, _, _, _, mockUpdateUseCase, _ := setup()
			id, err := domain.NewTodoID()
			require.NoError(t, err)
			todoID, _ := domain.NewTodoIDFromString(id.String())
			command := update.Command{
				ID:          *todoID,
				Title:       "test",
				Description: "description",
			}
			mockUpdateUseCase.On("Execute", mock.Anything, command).Return(fmt.Errorf("usecase error"))

			req := httptest.NewRequest(http.MethodPut, "/todos/"+id.String(), strings.NewReader(`{"title": "test", "description": "description"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusInternalServerError, rec.Code)
			mockUpdateUseCase.AssertExpectations(t)
		})

		t.Run("データが存在しない場合", func(t *testing.T) {
			e, _, _, _, _, mockUpdateUseCase, _ := setup()
			id, err := domain.NewTodoID()
			require.NoError(t, err)
			todoID, _ := domain.NewTodoIDFromString(id.String())
			command := update.Command{
				ID:          *todoID,
				Title:       "test",
				Description: "description",
			}
			mockUpdateUseCase.On("Execute", mock.Anything, command).Return(&domainError.ErrorTodoNotFound{ID: id.String()})

			req := httptest.NewRequest(http.MethodPut, "/todos/"+id.String(), strings.NewReader(`{"title": "test", "description": "description"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusNotFound, rec.Code)
			require.JSONEq(t, `{"message":"Todo not found"}`, rec.Body.String())
			mockUpdateUseCase.AssertExpectations(t)
		})
	})

	t.Run("DeleteTodo", func(t *testing.T) {
		t.Run("正常系", func(t *testing.T) {
			e, _, _, _, _, _, mockDeleteUseCase := setup()
			id, err := domain.NewTodoID()
			require.NoError(t, err)
			todoID, _ := domain.NewTodoIDFromString(id.String())
			command := delete.Command{
				ID: *todoID,
			}
			mockDeleteUseCase.On("Execute", mock.Anything, command).Return(nil)

			req := httptest.NewRequest(http.MethodDelete, "/todos/"+id.String(), nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusNoContent, rec.Code)
			require.Equal(t, "", rec.Body.String())
			mockDeleteUseCase.AssertExpectations(t)
		})

		t.Run("IDがUUIDじゃない場合", func(t *testing.T) {
			e, _, _, _, _, _, _ := setup()
			req := httptest.NewRequest(http.MethodDelete, "/todos/invalid-uuid", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
			require.JSONEq(t, `{"message":"Invalid todo ID"}`, rec.Body.String())
		})

		t.Run("ユースケースでエラーが発生したパターン", func(t *testing.T) {
			e, _, _, _, _, _, mockDeleteUseCase := setup()
			id, err := domain.NewTodoID()
			require.NoError(t, err)
			todoID, _ := domain.NewTodoIDFromString(id.String())
			command := delete.Command{
				ID: *todoID,
			}
			mockDeleteUseCase.On("Execute", mock.Anything, command).Return(fmt.Errorf("usecase error"))

			req := httptest.NewRequest(http.MethodDelete, "/todos/"+id.String(), nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusInternalServerError, rec.Code)
			mockDeleteUseCase.AssertExpectations(t)
		})

		t.Run("データが存在しない場合", func(t *testing.T) {
			e, _, _, _, _, _, mockDeleteUseCase := setup()
			id, err := domain.NewTodoID()
			require.NoError(t, err)
			todoID, _ := domain.NewTodoIDFromString(id.String())
			command := delete.Command{
				ID: *todoID,
			}
			mockDeleteUseCase.On("Execute", mock.Anything, command).Return(&domainError.ErrorTodoNotFound{ID: id.String()})

			req := httptest.NewRequest(http.MethodDelete, "/todos/"+id.String(), nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			require.Equal(t, http.StatusNotFound, rec.Code)
			require.JSONEq(t, `{"message":"Todo not found"}`, rec.Body.String())
			mockDeleteUseCase.AssertExpectations(t)
		})
	})
}
