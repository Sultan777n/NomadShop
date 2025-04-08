package models

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID         uint        `gorm:"primaryKey"`
	UserID     uint        `gorm:"not null"`
	OrderDate  time.Time   `gorm:"not null"`
	Status     string      `gorm:"not null"` // "pending", "completed", "shipped", т.б.
	Total      float64     `gorm:"not null"`
	User       User        `gorm:"foreignKey:UserID;references:ID"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;references:ID"`
}

// Миграция функциясы
func CreateOrder(db *gorm.DB, order *Order) (*Order, error) {
	err := db.Create(&order).Error
	return order, err
}

func GetOrdersByUser(db *gorm.DB, userID uint) ([]Order, error) {
	var orders []Order
	err := db.Where("user_id = ?", userID).Preload("OrderItems").Preload("User").Find(&orders).Error
	return orders, err
}

func GetOrderByID(db *gorm.DB, orderID uint) (*Order, error) {
	var order Order
	err := db.Where("id = ?", orderID).Preload("OrderItems").Preload("User").First(&order).Error
	return &order, err
}
