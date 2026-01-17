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
