package user

import (
	"moapick/db"
	"moapick/db/models"
)

func GetUserById(userId uint) (models.User, error) {
	var user models.User

	result := db.Client.First(&user, "id = ?", userId)

	return user, result.Error
}
