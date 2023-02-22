package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/go-api-fiber/controllers"
	"github.com/mr-emerald-wolf/go-api-fiber/middleware"
)

func UserRoutes(app *fiber.App) {
	user := app.Group("/user", middleware.DeserializeUser)
	user.Post("/create", controllers.CreateUser)
	user.Get("/findall", controllers.FindUser)
	user.Delete("/delete/:userId", controllers.DeleteUser)
	user.Get("/find/:userId", controllers.FindUserById)
	user.Put("/update/:userId", controllers.UpdateUser)
	user.Get("/me", controllers.GetMe)

}
