package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&Todo{})
	require.NoError(t, err)

	return db
}

func TestTodoRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTodoRepository(db)
	ctx := context.Background()

	todos := []Todo{
		{Title: "test title1", Description: "test description1"},
		{Title: "test title2", Description: "test description2"},
		{Title: "test title3", Description: "test description3"},
	}
	result := db.CreateInBatches(todos, 100)
	require.NoError(t, result.Error)

	allTodos, err := repo.FindAll(ctx)
	require.NoError(t, err)

	require.Equal(t, 3, len(allTodos))
	for i := range allTodos {
		require.Equal(t, todos[i].Title, allTodos[i].Title)
		require.Equal(t, todos[i].Description, allTodos[i].Description)
	}
}

func TestTodoRepository_FindAll_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTodoRepository(db)
	ctx := context.Background()

	allTodos, err := repo.FindAll(ctx)
	require.NoError(t, err)

	require.Equal(t, 0, len(allTodos))
}

func TestTodoRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTodoRepository(db)
	ctx := context.Background()

	todo := Todo{Title: "test title1", Description: "test description1"}
	result := db.Create(&todo)
	require.NoError(t, result.Error)

	todoResult, err := repo.FindByID(ctx, todo.ID)
	require.NoError(t, err)

	require.Equal(t, todo.Title, todoResult.Title)
	require.Equal(t, todo.Description, todoResult.Description)
}

func TestTodoRepository_FindByID_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTodoRepository(db)
	ctx := context.Background()

	_, err := repo.FindByID(ctx, 1)
	require.Error(t, err)
}
