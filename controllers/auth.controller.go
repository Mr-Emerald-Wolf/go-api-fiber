package controllers

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/go-api-fiber/initializers"
	"github.com/mr-emerald-wolf/go-api-fiber/models"
	"github.com/mr-emerald-wolf/go-api-fiber/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	randomState       = "pseudo-random"
)

func init() {
	initializers.LoadEnvVariables()

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

}

func HandleLogin(c *fiber.Ctx) error {
	url := googleOauthConfig.AuthCodeURL(randomState)
	c.Redirect(url)
	return nil
}

func HandleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	code := c.Query("code")

	if state != randomState {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "State is not valid", "state": state})
	}
	if code == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Authorization code not provided!"})
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Could not get token"})
	}
	googleUser, err := utils.GetGoogleUser(token.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Could not get user"})
	}

	user := models.User{
		Name:      googleUser.Name,
		Email:     googleUser.Email,
	}
	// Find or create new user
	if initializers.DB.Model(&user).Where("email = ?", googleUser.Email).Updates(&user).RowsAffected == 0 {
		initializers.DB.Create(&user)
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": user})
}
