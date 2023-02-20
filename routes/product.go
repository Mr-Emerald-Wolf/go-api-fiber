package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/go-api-fiber/controllers"
)

func ProductRoutes(app *fiber.App) {
	product := app.Group("/product")
	product.Post("/create", controllers.CreateProduct)
	product.Get("/findall", controllers.FindProducts)
	product.Delete("/delete/:productId", controllers.DeleteProduct)
	product.Get("/find/:productId", controllers.FindProductById)
	product.Put("/update/:productId", controllers.UpdateProduct)

}
