package entity

import (
	"gorm.io/gorm"
	"time"
)

type Book struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null" json:"title"`
	Author      string         `gorm:"not null" json:"author"`
	Price       float64        `gorm:"not null" json:"price"`
	Stock       int            `gorm:"not null;default:0" json:"stock"`
	Year        int            `gorm:"not null" json:"year"`
	CategoryID  uint           `gorm:"not null" json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	ImageBase64 string         `gorm:"type:text" json:"image_base64,omitempty"`
}