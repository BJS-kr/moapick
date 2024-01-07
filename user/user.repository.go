package user

import (
	"moapick/db"
	"moapick/db/models"

	"gorm.io/gorm/clause"
)

func GetUserById(userId uint) (*models.User, error) {
	user := new(models.User)
	result := db.Client.First(user, "id = ?", userId)

	return user, result.Error
}

func GetUserByEmail(email string) (*models.User, error) {
	user := new(models.User)
	result := db.Client.First(user, "email = ?", email)

	return user, result.Error
}

func CreateUserIfNotExists(userEntity *models.User) error {
	result := db.Client.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).Create(userEntity)

	return result.Error
}
