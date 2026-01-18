package repository

import (
	"context"
	"echo-sample2/internal/todo/domain"
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

		err := db.Create(todo1).Error
		require.NoError(t, err)
		err = db.Create(todo2).Error
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
}
