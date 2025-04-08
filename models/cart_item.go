package models

import (
	"gorm.io/gorm"
)

type CartItem struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	Quantity  uint    `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID"`
}

func AddToCart(db *gorm.DB, cartItem *CartItem) (*CartItem, error) {
	err := db.Create(&cartItem).Error
	if err != nil {
		return nil, err
	}

	err = db.Preload("Product").Preload("Product.Category").First(&cartItem, cartItem.ID).Error
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}

func GetCartItems(db *gorm.DB, userID uint) ([]CartItem, error) {
	var cartItems []CartItem
	err := db.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error
	return cartItems, err
}

func GetCartItemByID(db *gorm.DB, id uint) (*CartItem, error) {
	var cartItem CartItem
	err := db.Preload("Product").First(&cartItem, id).Error
	if err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func GetCartItemsByProductID(db *gorm.DB, productID uint) ([]CartItem, error) {
	var cartItems []CartItem
	// `product_id` арқылы іздеу және байланысты өнімнің мәліметтерін жүктеу
	err := db.Preload("Product").Where("product_id = ?", productID).Find(&cartItems).Error
	return cartItems, err
}
