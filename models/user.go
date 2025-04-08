package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"not null;unique"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}

func CreateUser(db *gorm.DB, user *User) (*User, error) {
	err := db.Create(&user).Error
	return user, err
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	err := db.First(&user, id).Error
	return &user, err
}

func GetUsers(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}

func UpdateUser(db *gorm.DB, id uint, updatedUser *User) (*User, error) {
	var user User
	err := db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	// Қолданушыны жаңарту
	err = db.Model(&user).Updates(updatedUser).Error
	return &user, err
}

func DeleteUser(db *gorm.DB, id uint) error {
	err := db.Delete(&User{}, id).Error
	return err
}
