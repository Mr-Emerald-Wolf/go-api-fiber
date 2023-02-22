package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/go-api-fiber/controllers"
)

func AuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Get("/login", controllers.HandleLogin)
	auth.Get("/callback", controllers.HandleCallback)

}
