package user

import (
	"errors"
	"log"
	"moapick/db/models"

	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SignInBody struct {
	Email string `json:"email" binding:"required"`
}
type JwtAccessToken struct {
	AccessToken string `json:"access_token"`
}

type UserController struct {
	UserRepository
	UserService
}

func (uc UserController)SignIn(c *fiber.Ctx) error {
	singInBody := new(SignInBody)

	if err := c.BodyParser(singInBody); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON( fiber.Map{"error": err.Error()})
	}

	newUser := models.User{Email: singInBody.Email}
	err := uc.UserRepository.CreateUserIfNotExists(&newUser)

	if err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	user, err := uc.UserRepository.GetUserByEmail(singInBody.Email)

	if err != nil {
		log.Println(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if jwt, err := uc.UserService.IssueJwt(singInBody.Email, user.ID); err == nil {
		responseBody := JwtAccessToken{AccessToken: jwt}
		return c.JSON(responseBody)
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "error during sign-in process",
	})
}

func (uc UserController)GetUserById(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("userId"))
	
	if err != nil {
		return c.JSON(fiber.StatusBadRequest, "userId must be integer")
	}

	user, err := uc.UserRepository.GetUserById(uint(userId))

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusBadRequest, "user not found by given id")
		} else {
			return c.JSON(http.StatusInternalServerError, "Internal Server Error")
		}
	}

	return c.JSON(user)
}

