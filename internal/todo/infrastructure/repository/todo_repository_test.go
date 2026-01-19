package repository

import (
	"context"
	"echo-sample2/internal/todo/domain"
	domainErr "echo-sample2/internal/todo/domain/errors"
	"echo-sample2/internal/todo/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTodoRepository_Save(t *testing.T) {
	// Setup
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&infrastructure.TodoEntity{})
	require.NoError(t, err)

	repo := NewTodoRepository(db)
	ctx := context.Background()

	t.Run("正常に保存できること", func(t *testing.T) {
		todo := &domain.Todo{
			ID:          domain.NewTodoID(),
			Title:       "Test Title",
			Description: "Test Description",
		}

		err := repo.Save(ctx, todo)
		assert.NoError(t, err)

		// 検証
		var entity infrastructure.TodoEntity
		err = db.First(&entity, "id = ?", todo.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, todo.ID, entity.ID)
		assert.Equal(t, todo.Title, entity.Title)
		assert.Equal(t, todo.Description, entity.Description)
	})

	t.Run("データベースエラーが発生した場合、エラーを返すこと", func(t *testing.T) {
		// テーブルをドロップしてエラーを発生させる
		err := db.Migrator().DropTable(&infrastructure.TodoEntity{})
		require.NoError(t, err)

		todo := &domain.Todo{
			ID:          domain.NewTodoID(),
			Title:       "Test Title",
			Description: "Test Description",
		}

		err = repo.Save(ctx, todo)
		assert.Error(t, err)

		// 次のテストのためにテーブルを再作成
		err = db.AutoMigrate(&infrastructure.TodoEntity{})
		require.NoError(t, err)
	})
}

func TestTodoRepository_FindAll(t *testing.T) {
	// Setup
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&infrastructure.TodoEntity{})
	require.NoError(t, err)

	repo := NewTodoRepository(db)
	ctx := context.Background()

	t.Run("全てのTodoを取得できること", func(t *testing.T) {
		todo1 := &infrastructure.TodoEntity{
			ID:          domain.NewTodoID(),
			Title:       "Title 1",
			Description: "Description 1",
		}
		todo2 := &infrastructure.TodoEntity{
			ID:          domain.NewTodoID(),
			Title:       "Title 2",
			Description: "Description 2",
		}

		err := db.Create([]*infrastructure.TodoEntity{todo1, todo2}).Error
		require.NoError(t, err)

		todos, err := repo.FindAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, todos, 2)

		assert.Equal(t, todo1.ID, todos[0].ID)
		assert.Equal(t, todo1.Title, todos[0].Title)
		assert.Equal(t, todo2.ID, todos[1].ID)
		assert.Equal(t, todo2.Title, todos[1].Title)
	})

	t.Run("データがない場合は空のリストを返すこと", func(t *testing.T) {
		// 新しいDB（メモリ内）を作成して空の状態にする
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		require.NoError(t, err)
		err = db.AutoMigrate(&infrastructure.TodoEntity{})
		require.NoError(t, err)
		repo := NewTodoRepository(db)

		todos, err := repo.FindAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, todos, 0)
	})

	t.Run("データベースエラーが発生した場合、エラーを返すこと", func(t *testing.T) {
		// テーブルをドロップしてエラーを発生させる
		err := db.Migrator().DropTable(&infrastructure.TodoEntity{})
		require.NoError(t, err)

		todos, err := repo.FindAll(ctx)
		assert.Error(t, err)
		assert.Nil(t, todos)
	})
}

func TestTodoRepository_FindById(t *testing.T) {
	// Setup
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&infrastructure.TodoEntity{})
	require.NoError(t, err)

	repo := NewTodoRepository(db)
	ctx := context.Background()

	t.Run("IDを指定してTodoを取得できること", func(t *testing.T) {
		todo := &infrastructure.TodoEntity{
			ID:          domain.NewTodoID(),
			Title:       "Title",
			Description: "Description",
		}

		err := db.Create(todo).Error
		require.NoError(t, err)

		found, err := repo.FindById(ctx, todo.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, todo.ID, found.ID)
		assert.Equal(t, todo.Title, found.Title)
		assert.Equal(t, todo.Description, found.Description)
	})

	t.Run("存在しないIDの場合はnilを返すこと", func(t *testing.T) {
		found, err := repo.FindById(ctx, domain.NewTodoID())
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("データベースエラーが発生した場合、エラーを返すこと", func(t *testing.T) {
		// テーブルをドロップしてエラーを発生させる
		err := db.Migrator().DropTable(&infrastructure.TodoEntity{})
		require.NoError(t, err)

		found, err := repo.FindById(ctx, domain.NewTodoID())
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestTodoRepository_Update(t *testing.T) {
	// Setup
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&infrastructure.TodoEntity{})
	require.NoError(t, err)

	repo := NewTodoRepository(db)
	ctx := context.Background()

	t.Run("正常に更新できること", func(t *testing.T) {
		id := domain.NewTodoID()
		initialTodo := &infrastructure.TodoEntity{
			ID:          id,
			Title:       "Initial Title",
			Description: "Initial Description",
		}
		err := db.Create(initialTodo).Error
		require.NoError(t, err)

		updatedTodo := &domain.Todo{
			ID:          id,
			Title:       "Updated Title",
			Description: "Updated Description",
		}

		err = repo.Update(ctx, updatedTodo)
		assert.NoError(t, err)

		// 検証
		var entity infrastructure.TodoEntity
		err = db.First(&entity, "id = ?", id).Error
		assert.NoError(t, err)
		assert.Equal(t, updatedTodo.Title, entity.Title)
		assert.Equal(t, updatedTodo.Description, entity.Description)
	})

	t.Run("存在しないIDの場合はErrorTodoNotFoundを返すこと", func(t *testing.T) {
		todo := &domain.Todo{
			ID:          domain.NewTodoID(),
			Title:       "Title",
			Description: "Description",
		}

		err := repo.Update(ctx, todo)
		assert.Error(t, err)
		assert.IsType(t, &domainErr.ErrorTodoNotFound{}, err)
	})

	t.Run("データベースエラーが発生した場合、エラーを返すこと", func(t *testing.T) {
		// テーブルをドロップしてエラーを発生させる
		err := db.Migrator().DropTable(&infrastructure.TodoEntity{})
		require.NoError(t, err)

		todo := &domain.Todo{
			ID:          domain.NewTodoID(),
			Title:       "Title",
			Description: "Description",
		}

		err = repo.Update(ctx, todo)
		assert.Error(t, err)
	})
}

func TestTodoRepository_Delete(t *testing.T) {
	// Setup
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&infrastructure.TodoEntity{})
	require.NoError(t, err)

	repo := NewTodoRepository(db)
	ctx := context.Background()

	t.Run("正常に削除できること", func(t *testing.T) {
		id := domain.NewTodoID()
		todo := &infrastructure.TodoEntity{
			ID:          id,
			Title:       "Test Title",
			Description: "Test Description",
		}
		err := db.Create(todo).Error
		require.NoError(t, err)

		err = repo.Delete(ctx, id)
		assert.NoError(t, err)

		// 検証: 削除されていること
		var entity infrastructure.TodoEntity
		err = db.First(&entity, "id = ?", id).Error
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("存在しないIDの場合はErrorTodoNotFoundを返すこと", func(t *testing.T) {
		err := repo.Delete(ctx, domain.NewTodoID())
		assert.Error(t, err)
		assert.IsType(t, &domainErr.ErrorTodoNotFound{}, err)
	})

	t.Run("データベースエラーが発生した場合、エラーを返すこと", func(t *testing.T) {
		// テーブルをドロップしてエラーを発生させる
		err := db.Migrator().DropTable(&infrastructure.TodoEntity{})
		require.NoError(t, err)

		err = repo.Delete(ctx, domain.NewTodoID())
		assert.Error(t, err)

		// 次のテストのために（もしあれば）テーブルを再作成
		err = db.AutoMigrate(&infrastructure.TodoEntity{})
		require.NoError(t, err)
	})
}
