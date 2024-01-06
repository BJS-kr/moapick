package user

import (
	"errors"
	"fmt"
	"log"
	"moapick/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SignIn struct {
	Email string `json:"email" binding:"required"`
}
type JwtAccessToken struct {
	AccessToken string `json:"access_token"`
}

func UserController(r *gin.Engine) {
	u := r.Group("/user")

	u.POST("/sign-in", func(c *gin.Context) {
		singInBody := SignIn{}

		if err := c.ShouldBind(&singInBody); err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if jwt, err := IssueJwt(singInBody.Email); err == nil {
			responseBody := JwtAccessToken{AccessToken: jwt}
			c.JSON(http.StatusCreated, responseBody)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error during sign-in process",
			})
		}
	})

	au := u.Group("/")
	au.Use(middleware.JwtMiddleware())

	au.GET("/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		user, err := GetUserById(userId)

		if err != nil {
			handleFindOneError(c, err, "User", "userId")
		} else {
			c.JSON(http.StatusOK, user)
		}
	})
}

func handleFindOneError(c *gin.Context, err error, target, by string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s not found by given %s", target, by))
	} else {
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	log.Print(err.Error())
}
