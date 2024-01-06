package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserController(r *gin.Engine) {
	userCont := r.Group("/user")

	userCont.POST("/sing-in", func(ctx *gin.Context) {

	})
	userCont.DELETE("/sign-out", func(ctx *gin.Context) {})

	userCont.GET("/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		user, err := GetUser(userId)

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
}
