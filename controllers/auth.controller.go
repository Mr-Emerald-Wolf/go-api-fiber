package controllers

import (
	"context"
	"os"
	"time"

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

	googleToken, err := googleOauthConfig.Exchange(context.Background(), code)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Could not get token"})
	}
	googleUser, err := utils.GetGoogleUser(googleToken.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Could not get user"})
	}

	user := models.User{
		Name:  googleUser.Name,
		Email: googleUser.Email,
	}
	// Find or create new user
	if initializers.DB.Model(&user).Where("email = ?", googleUser.Email).Updates(&user).RowsAffected == 0 {
		initializers.DB.Create(&user)
	}

	duration, _ := time.ParseDuration("1h")
	token, err := utils.GenerateToken(duration, user.Email, os.Getenv("JWTTokenSecret"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Could not create JWT"})
	}

	// Create cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.HTTPOnly = true
	cookie.Secure = false
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"

	// Set cookie
	c.Cookie(cookie)

	return c.Redirect("http://localhost:3000")
}
