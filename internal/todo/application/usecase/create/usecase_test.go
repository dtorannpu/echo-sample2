package create

import (
	"context"
	"echo-sample2/internal/todo/domain"
	"echo-sample2/internal/todo/domain/repository"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// モックの定義
type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) Do(ctx context.Context, fn func(tx repository.Transaction) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	// fnを呼び出して、モックのTransactionを渡す
	return fn(args.Get(1).(repository.Transaction))
}

type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) TodoRepo() repository.TodoRepository {
	args := m.Called()
	return args.Get(0).(repository.TodoRepository)
}

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) Save(ctx context.Context, todo *domain.Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func TestCreateTodoUseCase_Execute(t *testing.T) {
	t.Run("正常にTodoが作成されること", func(t *testing.T) {
		// 準備
		mockTM := new(MockTransactionManager)
		mockTX := new(MockTransaction)
		mockRepo := new(MockTodoRepository)

		uc := New(mockTM)

		ctx := context.Background()
		command := Command{
			Title:       "テストタイトル",
			Description: "テスト説明",
		}

		// モックの挙動を設定
		mockTM.On("Do", ctx, mock.Anything).Return(nil, mockTX)
		mockTX.On("TodoRepo").Return(mockRepo)
		mockRepo.On("Save", ctx, mock.MatchedBy(func(todo *domain.Todo) bool {
			return todo.Title == command.Title && todo.Description == command.Description && todo.ID != domain.TodoID{}
		})).Return(nil)

		// 実行
		err := uc.Execute(ctx, command)

		// 検証
		require.NoError(t, err)
		mockTM.AssertExpectations(t)
		mockTX.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("リポジトリの保存でエラーが発生した場合、エラーを返すこと", func(t *testing.T) {
		// 準備
		mockTM := new(MockTransactionManager)
		mockTX := new(MockTransaction)
		mockRepo := new(MockTodoRepository)

		uc := New(mockTM)

		ctx := context.Background()
		command := Command{
			Title:       "テストタイトル",
			Description: "テスト説明",
		}

		expectedErr := context.DeadlineExceeded

		// モックの挙動を設定
		mockTM.On("Do", ctx, mock.Anything).Return(nil, mockTX)
		mockTX.On("TodoRepo").Return(mockRepo)
		mockRepo.On("Save", ctx, mock.Anything).Return(expectedErr)

		// 実行
		err := uc.Execute(ctx, command)

		// 検証
		require.ErrorIs(t, err, expectedErr)
		mockTM.AssertExpectations(t)
		mockTX.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}
