package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const BearerSchema = "Bearer "

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, BearerSchema)

		parserOption := jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name})
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, assertionSuccess := t.Method.(*jwt.SigningMethodHMAC); assertionSuccess {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(os.Getenv("SECRET_KEY")), nil
		}, parserOption)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unexpected access token algorithm"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			email, found := claims["email"]
			if !found {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "expected claim missing"})
			}

			c.Set("email", email)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unexpected claims found"})
			c.Abort()
			return
		}

		c.Next()
	}
}
