package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"not null;unique"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`

	Roles []Role `gorm:"many2many:user_roles;" json:"roles"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.Password != "" && !isBcryptHashed(u.Password) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

func isBcryptHashed(pw string) bool {
	return len(pw) == 60 && (pw[:4] == "$2a$" || pw[:4] == "$2b$" || pw[:4] == "$2y$")
}

func CreateUser(db *gorm.DB, user *User) (*User, error) {
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	err := db.Preload("Roles").First(&user, id).Error
	return &user, err
}

func GetUsers(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Preload("Roles").Find(&users).Error
	return users, err
}

func UpdateUser(db *gorm.DB, id uint, updatedUser *User) (*User, error) {
	var user User
	err := db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	err = db.Model(&user).Updates(updatedUser).Error
	return &user, err
}

func DeleteUser(db *gorm.DB, id uint) error {
	if err := db.Where("user_id = ?", id).Delete(&UserRole{}).Error; err != nil {
		return err
	}

	return db.Delete(&User{}, id).Error
}
