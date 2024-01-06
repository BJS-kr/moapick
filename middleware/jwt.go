package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const BearerSchema = "Bearer "

func JwtMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error{
		authHeader := c.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, BearerSchema)

		parserOption := jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name})
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, assertionSuccess := t.Method.(*jwt.SigningMethodHMAC); assertionSuccess {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(os.Getenv("SECRET_KEY")), nil
		}, parserOption)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON( fiber.Map{"error": "unexpected access token algorithm"})
			
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			email, found := claims["email"]

			if !found {
				return c.Status(fiber.StatusUnauthorized).JSON( fiber.Map{"error": "expected claim missing"})
			
			}

			c.Locals("email", email)
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON( fiber.Map{"error": "unexpected claims found"})
		
		}

		return c.Next()
	}
}
