package repository

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title       string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:varchar(1000);not null"`
}
