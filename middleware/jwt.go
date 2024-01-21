package middleware

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const BearerSchema = "Bearer "

func JwtMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, BearerSchema)

		parserOption := jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name})
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, assertionSuccess := t.Method.(*jwt.SigningMethodHMAC); !assertionSuccess {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(os.Getenv("SECRET_KEY")), nil
		}, parserOption)

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unexpected access token algorithm"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			email, emailFound := claims["email"]
			userId, userIdFound := claims["user_id"]

			if !emailFound || !userIdFound {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "expected claim missing"})
			}

			floatUserId, ok := userId.(float64)

			if !ok {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
			}

			c.Locals("email", email)
			c.Locals("userId", uint(floatUserId))
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unexpected claims found"})
		}

		return c.Next()
	}
}
