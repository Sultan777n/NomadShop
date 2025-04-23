package models

import (
	"gorm.io/gorm"
)

type Product struct {
	ID          uint     `gorm:"primaryKey" json:"id"`
	Name        string   `json:"name"`
	Price       uint     `json:"price"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Color       string   `json:"color"`
	Size        string   `json:"size"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `gorm:"foreignKey:CategoryID" json:"category"`
	Stock       uint     `json:"stock"`
}

func GetProducts(db *gorm.DB) ([]Product, error) {
	var products []Product
	err := db.Preload("Category").Find(&products).Error
	return products, err
}

func CreateProduct(db *gorm.DB, product *Product) (*Product, error) {
	err := db.Create(product).Error
	return product, err
}

func GetProductByID(db *gorm.DB, id uint) (*Product, error) {
	var product Product
	err := db.Preload("Category").First(&product, id).Error
	return &product, err
}

func UpdateProduct(db *gorm.DB, id uint, updated *Product) (*Product, error) {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return nil, err
	}

	err := db.Model(&product).Updates(updated).Error
	if err != nil {
		return nil, err
	}

	// Категорияны қайта жүктеу
	db.Preload("Category").First(&product, id)
	return &product, nil
}

func DeleteProduct(db *gorm.DB, id uint) error {
	return db.Delete(&Product{}, id).Error
}
