package user

import (
	"moapick/db"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
}

func GetUserById(userId string) (User, error) {
	var user User

	result := db.DB.First(&user, "id = ?", userId)

	return user, result.Error
}
