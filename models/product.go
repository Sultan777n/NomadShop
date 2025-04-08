package models

import (
	"gorm.io/gorm"
)

type Product struct {
	ID          uint     `gorm:"primaryKey"`
	Name        string   `gorm:"not null"`
	Price       uint     `gorm:"not null"`
	Description string   `gorm:"not null"`
	Image       string   `gorm:"not null"`
	Color       string   `gorm:"not null"`
	Size        string   `gorm:"not null"`
	CategoryID  uint     `gorm:"not null"`
	Category    Category `gorm:"foreignKey:CategoryID"`
	Stock       uint     `gorm:"not null"`
}

func GetProducts(db *gorm.DB) ([]Product, error) {
	var products []Product
	err := db.Preload("Category").Find(&products).Error
	return products, err
}

func CreateProduct(db *gorm.DB, product *Product) (*Product, error) {
	err := db.Create(&product).Error
	return product, err
}

func GetProductByID(db *gorm.DB, id uint) (*Product, error) {
	var product Product
	err := db.Preload("Category").First(&product, id).Error
	return &product, err
}

func UpdateProduct(db *gorm.DB, id uint, product *Product) (*Product, error) {
	err := db.Model(&Product{}).Where("id = ?", id).Updates(product).Error
	return product, err
}

func DeleteProduct(db *gorm.DB, id uint) error {
	err := db.Delete(&Product{}, id).Error
	return err
}
