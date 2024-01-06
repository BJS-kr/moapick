package user

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)



func IssueJwt(userEmail string) (string, error) {
	var (
		key   []byte
		token *jwt.Token
	)

	key = []byte(os.Getenv("SECRET_KEY"))
	token = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":   "moapick",
			"email": userEmail,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	return token.SignedString(key)
}
