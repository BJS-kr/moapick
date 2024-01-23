package user

import (
	"errors"
	"log"
	"moapick/common"
	"moapick/db/models"

	"net/http"

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

// ShowAccount godoc
//	@Summary		sign in
//	@Description	미리 가입되어있어야 하거나 비밀번호 같은 것 필요없습니다. 그냥 이메일만 보내면 그에 맞는 토큰을 생성해 돌려줍니다.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SignInBody	true	"email to login"	Format(email)
//	@Success		201		{object}	JwtAccessToken
//	@Failure		400		{object}	common.ErrorMessage
//	@Failure		500		{object}	common.ErrorMessage
//	@Router			/user/sign-in [post]
func (uc UserController)SignIn(c *fiber.Ctx) error {
	singInBody := new(SignInBody)

	if err := c.BodyParser(singInBody); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorMessage{Error: err.Error()})
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
		return c.Status(fiber.StatusCreated).JSON(responseBody)
	}

	return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{
		Error: "error during sign-in process",
	})
}
// ShowAccount godoc
//	@Summary		get user
//	@Description	user id를 통해 유저 정보를 반환합니다.
//	@Tags			user
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	models.User
//	@Failure		400				{object}	common.ErrorMessage
//	@Failure		500				{object}	common.ErrorMessage
//	@Router			/user [get]
func (uc UserController)GetUserById(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(uint)
	
	if !ok {
		return c.JSON(fiber.StatusBadRequest, "failed to get user id")
	}

	user, err := uc.UserRepository.GetUserById(userId)

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

