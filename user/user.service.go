package user

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserService struct{}

func (us UserService) IssueJwt(userEmail string, userId uint) (string, error) {
	var (
		key   []byte
		token *jwt.Token
	)

	key = []byte(os.Getenv("SECRET_KEY"))
	token = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":     "moapick",
			"email":   userEmail,
			"user_id": userId,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

	return token.SignedString(key)
}
