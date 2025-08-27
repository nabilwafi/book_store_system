package entity

import (
	"gorm.io/gorm"
	"time"
)

type Category struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"uniqueIndex;not null" json:"name"`
	Books     []Book         `gorm:"foreignKey:CategoryID" json:"books,omitempty"`
}