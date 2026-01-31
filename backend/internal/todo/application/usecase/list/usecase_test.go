package list

import (
	"context"
	"echo-sample2/internal/todo/domain"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockGetTodosRepository struct {
	mock.Mock
}

func (m *MockGetTodosRepository) FindAll(ctx context.Context) ([]*domain.Todo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Todo), args.Error(1)
}

func TestGetTodosUseCase_Execute(t *testing.T) {
	t.Run("正常にTodo一覧が取得できること", func(t *testing.T) {
		// 準備
		mockRepo := new(MockGetTodosRepository)
		uc := New(mockRepo)
		ctx := context.Background()

		id1 := domain.NewTodoID()
		id2 := domain.NewTodoID()

		todos := []*domain.Todo{
			{
				ID:          id1,
				Title:       "Title 1",
				Description: "Description 1",
			},
			{
				ID:          id2,
				Title:       "Title 2",
				Description: "Description 2",
			},
		}

		mockRepo.On("FindAll", ctx).Return(todos, nil)

		// 実行
		result, err := uc.Execute(ctx)

		// 検証
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Todos, 2)

		require.Equal(t, id1, result.Todos[0].ID)
		require.Equal(t, "Title 1", result.Todos[0].Title)
		require.Equal(t, "Description 1", result.Todos[0].Description)

		require.Equal(t, id2, result.Todos[1].ID)
		require.Equal(t, "Title 2", result.Todos[1].Title)
		require.Equal(t, "Description 2", result.Todos[1].Description)

		mockRepo.AssertExpectations(t)
	})

	t.Run("リポジトリでエラーが発生した場合、エラーを返すこと", func(t *testing.T) {
		// 準備
		mockRepo := new(MockGetTodosRepository)
		uc := New(mockRepo)
		ctx := context.Background()

		expectedErr := errors.New("repository error")
		mockRepo.On("FindAll", ctx).Return(nil, expectedErr)

		// 実行
		result, err := uc.Execute(ctx)

		// 検証
		require.ErrorIs(t, err, expectedErr)
		require.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Todoが空の場合、空のリストを返すこと", func(t *testing.T) {
		// 準備
		mockRepo := new(MockGetTodosRepository)
		uc := New(mockRepo)
		ctx := context.Background()

		mockRepo.On("FindAll", ctx).Return([]*domain.Todo{}, nil)

		// 実行
		result, err := uc.Execute(ctx)

		// 検証
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Empty(t, result.Todos)

		mockRepo.AssertExpectations(t)
	})
}
