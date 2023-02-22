package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/go-api-fiber/initializers"
	"github.com/mr-emerald-wolf/go-api-fiber/models"
	"github.com/mr-emerald-wolf/go-api-fiber/utils"
)

func init() {
	initializers.LoadEnvVariables()
}
func DeserializeUser(c *fiber.Ctx) error {

	var token string
	cookie := c.Cookies("token")

	// authorizationHeader := c.Request.Header.Get("Authorization")
	// fields := strings.Fields(authorizationHeader)

	// if len(fields) != 0 && fields[0] == "Bearer" {
	// 	token = fields[1]
	// } else if err == nil {
	// 	token = cookie
	// }

	token = cookie

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	sub, err := utils.ValidateToken(token, os.Getenv("JWTTokenSecret"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var user models.User
	result := initializers.DB.First(&user, "email = ?", fmt.Sprint(sub))
	if result.Error != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no logger exists"})
	}

	c.Set("currentUser", user.Email)
	return c.Next()

}
