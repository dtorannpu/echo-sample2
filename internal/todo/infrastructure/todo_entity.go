package infrastructure

import (
	"echo-sample2/internal/todo/domain"
	"time"

	"gorm.io/gorm"
)

type TodoEntity struct {
	ID          domain.TodoID `gorm:"type:uuid;primaryKey"`
	Title       string        `gorm:"type:text;size:100;not null"`
	Description string        `gorm:"type:text;size:1000"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
