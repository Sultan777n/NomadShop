package models

import "gorm.io/gorm"

type OrderItem struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	Quantity  uint    `gorm:"not null"`
	Price     float64 `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID"`
}

func CreateOrderItem(db *gorm.DB, orderItem *OrderItem) (*OrderItem, error) {
	err := db.Create(&orderItem).Error
	return orderItem, err
}

func GetOrderItemsByOrderID(db *gorm.DB, orderID uint) ([]OrderItem, error) {
	var items []OrderItem
	err := db.Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}
