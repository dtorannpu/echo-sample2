package get

import (
	"context"
	"echo-sample2/internal/todo/domain"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) FindById(ctx context.Context, id domain.TodoID) (*domain.Todo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Todo), args.Error(1)
}

func TestUseCase_Execute(t *testing.T) {
	t.Run("正常にTodoが取得できること", func(t *testing.T) {
		// 準備
		mockRepo := new(MockTodoRepository)
		uc := New(mockRepo)
		ctx := context.Background()
		id, err := domain.NewTodoID()
		require.NoError(t, err)

		todo := &domain.Todo{
			ID:          id,
			Title:       "Test Title",
			Description: "Test Description",
		}

		mockRepo.On("FindById", ctx, id).Return(todo, nil)

		// 実行
		result, err := uc.Execute(ctx, Command{ID: id})

		// 検証
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, id, result.ID)
		require.Equal(t, "Test Title", result.Title)
		require.Equal(t, "Test Description", result.Description)

		mockRepo.AssertExpectations(t)
	})

	t.Run("リポジトリでエラーが発生した場合、エラーを返すこと", func(t *testing.T) {
		// 準備
		mockRepo := new(MockTodoRepository)
		uc := New(mockRepo)
		ctx := context.Background()
		id, err := domain.NewTodoID()
		require.NoError(t, err)

		expectedErr := errors.New("repository error")
		mockRepo.On("FindById", ctx, id).Return(nil, expectedErr)

		// 実行
		result, err := uc.Execute(ctx, Command{ID: id})

		// 検証
		require.ErrorIs(t, err, expectedErr)
		require.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("データが見つからない場合、nilを返すこと", func(t *testing.T) {
		// 準備
		mockRepo := new(MockTodoRepository)
		uc := New(mockRepo)
		ctx := context.Background()
		id, err := domain.NewTodoID()
		require.NoError(t, err)

		mockRepo.On("FindById", ctx, id).Return(nil, nil)

		// 実行
		result, err := uc.Execute(ctx, Command{ID: id})

		// 検証
		require.NoError(t, err)
		require.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})
}
