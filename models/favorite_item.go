package models

import (
	"gorm.io/gorm"
)

type FavoriteItem struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
}

func AddToFavorites(db *gorm.DB, favoriteItem *FavoriteItem) (*FavoriteItem, error) {
	err := db.Create(favoriteItem).Error
	return favoriteItem, err
}

func GetFavoriteItems(db *gorm.DB, userID uint) ([]FavoriteItem, error) {
	var favoriteItems []FavoriteItem
	err := db.Preload("Product").Preload("Product.Category").Where("user_id = ?", userID).Find(&favoriteItems).Error
	return favoriteItems, err
}

func GetFavoriteItemByID(db *gorm.DB, id uint) (*FavoriteItem, error) {
	var favoriteItem FavoriteItem
	err := db.Preload("Product").Preload("Product.Category").First(&favoriteItem, id).Error
	if err != nil {
		return nil, err
	}
	return &favoriteItem, nil
}
