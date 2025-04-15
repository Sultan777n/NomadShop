package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

func FixPasswords(db *gorm.DB) {
	var users []User
	db.Find(&users)

	for _, user := range users {
		if len(user.Password) < 60 {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Println("Hash error:", err)
				continue
			}

			user.Password = string(hashedPassword)
			db.Save(&user)
			log.Printf("Hashed: %s\n", user.Email)
		}
	}
}
