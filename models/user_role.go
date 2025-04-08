package models

import (
	"gorm.io/gorm"
)

type UserRole struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null;index"`
	RoleID uint `gorm:"not null;index"`
	User   User `gorm:"foreignKey:UserID;references:ID"`
	Role   Role `gorm:"foreignKey:RoleID;references:ID"`
}

func AddUserRole(db *gorm.DB, userRole *UserRole) (*UserRole, error) {
	// Егер рөл мен пайдаланушы дұрыс болса, қосу операциясын орындау
	err := db.Create(&userRole).Error
	return userRole, err
}

func GetUserRoles(db *gorm.DB, userID uint) ([]UserRole, error) {
	var userRoles []UserRole
	err := db.Where("user_id = ?", userID).Preload("Role").Preload("User").Find(&userRoles).Error
	return userRoles, err
}

func DeleteUserRole(db *gorm.DB, userID uint, roleID uint) error {
	// Пайдаланушының рөлін өшіру операциясы
	err := db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&UserRole{}).Error
	return err
}

func GetRoleByUserAndRoleID(db *gorm.DB, userID, roleID uint) (*UserRole, error) {
	var userRole UserRole
	err := db.Where("user_id = ? AND role_id = ?", userID, roleID).First(&userRole).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}
