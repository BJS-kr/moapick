package user

import (
	"errors"
	"fmt"
	"log"
	"moapick/db/models"
	"moapick/middleware"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

type SignIn struct {
	Email string `json:"email" binding:"required"`
}
type JwtAccessToken struct {
	AccessToken string `json:"access_token"`
}

func UserController(r *fiber.App) {
	u := r.Group("/user")

	u.Post("/sign-in", func(c *fiber.Ctx) error {
		singInBody := new(SignIn)

		if err := c.BodyParser(singInBody); err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON( fiber.Map{"error": err.Error()})
		}

		newUser := models.User{Email: singInBody.Email}
		err := CreateUserIfNotExists(&newUser)

		if err != nil {
			log.Println(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		user, err := GetUserByEmail(singInBody.Email)

		if err != nil {
			log.Println(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if jwt, err := IssueJwt(singInBody.Email, user.ID); err == nil {
			responseBody := JwtAccessToken{AccessToken: jwt}
			return c.JSON(responseBody)
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "error during sign-in process",
			})
		}
	})

	au := u.Group("/")
	au.Use(middleware.JwtMiddleware())

	au.Get("/:userId", func(c *fiber.Ctx) error {
		userId, err := strconv.Atoi(c.Params("userId"))
		
		if err != nil {
			return c.JSON(http.StatusBadRequest, "userId must be integer")
		
		}

		user, err := GetUserById(uint(userId))

		if err != nil {
			return handleFindOneError(c, err, "User", "userId")
		} else {
			return c.JSON(user)
		}
	})
}

func handleFindOneError(c *fiber.Ctx, err error, target, by string) error {
	log.Println(err.Error())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("%s not found by given %s", target, by))
	} else {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

}
