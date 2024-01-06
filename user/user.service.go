package user

import (
	"moapick/db/models"
)

func GetUser(userId string) (models.User, error) {
	return GetUserById(userId)
}
