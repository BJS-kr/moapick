package user

import (
	"moapick/db/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	Client *gorm.DB
}

func (ur UserRepository) GetUserById(userId uint) (*models.User, error) {
	user := new(models.User)
	result := ur.Client.First(user, "id = ?", userId)

	return user, result.Error
}

func (ur UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := new(models.User)
	result := ur.Client.First(user, "email = ?", email)

	return user, result.Error
}

func (ur UserRepository) CreateUserIfNotExists(userEntity *models.User) error {
	result := ur.Client.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).Create(userEntity)

	return result.Error
}
