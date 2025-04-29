package models

import (
	"gorm.io/gorm"
)

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null;unique"`
}

func GetRoles(db *gorm.DB) ([]Role, error) {
	var roles []Role
	err := db.Find(&roles).Error
	return roles, err
}

func CreateRole(db *gorm.DB, role *Role) (*Role, error) {
	err := db.Create(&role).Error
	return role, err
}

func GetRoleByID(db *gorm.DB, id uint) (*Role, error) {
	var role Role
	err := db.First(&role, id).Error
	return &role, err
}

func UpdateRole(db *gorm.DB, id uint, role *Role) (*Role, error) {
	err := db.Model(&Role{}).Where("id = ?", id).Updates(role).Error
	return role, err
}

func DeleteRole(db *gorm.DB, id uint) error {
	err := db.Delete(&Role{}, id).Error
	return err
}
