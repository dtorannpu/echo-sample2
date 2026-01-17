package db

import (
	"echo-sample2/internal/todo/infrastructure"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB() (*DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&infrastructure.TodoEntity{})
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
