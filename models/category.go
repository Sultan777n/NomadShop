package models

import (
	"gorm.io/gorm"
)

type Category struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
	URL  string `gorm:"not null"`
}

func GetAllCategories(db *gorm.DB) ([]Category, error) {
	var categories []Category
	err := db.Find(&categories).Error
	return categories, err
}

func GetCategoryByID(db *gorm.DB, id uint) (*Category, error) {
	var category Category
	err := db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func CreateCategory(db *gorm.DB, category *Category) (*Category, error) {
	err := db.Create(&category).Error
	return category, err
}
