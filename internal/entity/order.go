package entity

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uint           `gorm:"not null" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TotalPrice float64        `gorm:"not null" json:"total_price"`
	Status     string         `gorm:"not null;default:pending" json:"status"`
	PaymentURL string         `json:"payment_url,omitempty"`
	OrderItems []OrderItem    `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

type OrderItem struct {
	ID       uint    `gorm:"primarykey" json:"id"`
	OrderID  uint    `gorm:"not null" json:"order_id"`
	Order    Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	BookID   uint    `gorm:"not null" json:"book_id"`
	Book     Book    `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Quantity int     `gorm:"not null" json:"quantity"`
	Price    float64 `gorm:"not null" json:"price"`
}