package user

import (
	"moapick/db"
	"moapick/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(r *fiber.App) {
	userRepository := UserRepository{Client: db.Client}
	userService := UserService{}
	userController := UserController{UserRepository: userRepository, UserService: userService}

	u := r.Group("/user")
	u.Post("/sign-in", userController.SignIn)

	au := u.Group("/")
	au.Use(middleware.JwtMiddleware())

	au.Get("/:userId", userController.GetUserById)
}

