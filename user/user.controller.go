package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserController(r *gin.Engine) {
	r.GET("/user/:userId", func (c *gin.Context) {
		userId := c.Param("userId")
		user, error := GetUser(userId)

		if errors.Is(error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, "User not found by given ID")
		}

		c.JSON(http.StatusOK, user)
	})
}