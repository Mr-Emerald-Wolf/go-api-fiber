package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/go-api-fiber/initializers"
	"github.com/mr-emerald-wolf/go-api-fiber/models"
	"github.com/mr-emerald-wolf/go-api-fiber/utils"
	"gorm.io/gorm"
)

func CreateUser(c *fiber.Ctx) error {
	var payload *models.CreateUserSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	now := time.Now()
	newUser := models.User{
		Name:      payload.Name,
		Email:     payload.Email,
		Phone:     payload.Phone,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Email already exist, please use another email"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"note": newUser}})
}

func FindUser(c *fiber.Ctx) error {

	var users []models.User
	results := initializers.DB.Find(&users)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(users), "users": users})
}

func UpdateUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	var payload *models.UpdateUserSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var user models.User
	result := initializers.DB.First(&user, "id = ?", userId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No user with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.Name != "" {
		updates["name"] = payload.Name
	}
	if payload.Email != "" {
		updates["email"] = payload.Email
	}
	if payload.Phone != "" {
		updates["phone"] = payload.Phone
	}

	updates["updated_at"] = time.Now()

	fmt.Println(updates)

	initializers.DB.Model(&user).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}

func FindUserById(c *fiber.Ctx) error {
	userId := c.Params("userId")
	// Get the user
	user := models.User{}
	result := initializers.DB.First(&user, "id = ?", userId)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": user})
}

func DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("userId")
	// Get the user
	user := models.User{}
	result := initializers.DB.First(&user, "id = ?", userId)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "User Not Found", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}

	result = initializers.DB.Delete(&models.User{}, "id = ?", userId)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "User Not Deleted", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "User Deleted"})
}

func GetMe(c *fiber.Ctx) error {
	email := c.GetRespHeader("currentUser")
	user := models.User{}

	result := initializers.DB.First(&user, "email = ?", fmt.Sprint(email))
	if result.Error != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no logger exists"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": user})
}
