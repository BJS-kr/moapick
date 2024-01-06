package user

import (
	"moapick/db"
	"moapick/db/models"
)

func GetUserById(userId string) (models.User, error) {
	var user models.User

	result := db.Client.First(&user, "id = ?", userId)

	return user, result.Error
}
