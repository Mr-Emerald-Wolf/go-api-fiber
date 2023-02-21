package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/go-api-fiber/initializers"
	"github.com/mr-emerald-wolf/go-api-fiber/models"
	"github.com/mr-emerald-wolf/go-api-fiber/utils"
	"gorm.io/gorm"
)

func CreateProduct(c *fiber.Ctx) error {
	var payload *models.CreateProductSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	now := time.Now()
	newProduct := models.Product{
		Name:      payload.Name,
		Price:     payload.Price,
		Category:  payload.Category,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := initializers.DB.Create(&newProduct)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Title already exist, please use another title"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"note": newProduct}})
}

func FindProducts(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var products []models.Product
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&products)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(products), "products": products})
}

func UpdateProduct(c *fiber.Ctx) error {
	productId := c.Params("productId")

	var payload *models.UpdateProductSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var product models.Product
	result := initializers.DB.First(&product, "id = ?", productId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No product with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.Name != "" {
		updates["name"] = payload.Name
	}
	if payload.Category != "" {
		updates["category"] = payload.Category
	}
	if payload.Price != 0 {
		updates["price"] = payload.Price
	}

	if payload.Published != nil {
		updates["published"] = payload.Published
	}

	updates["updated_at"] = time.Now()

	fmt.Println(updates)

	initializers.DB.Model(&product).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"product": product}})
}

func FindProductById(c *fiber.Ctx) error {
	productId := c.Params("productId")
	// Get the product
	product := models.Product{}
	result := initializers.DB.First(&product, "id = ?", productId)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "product": product})
}

func DeleteProduct(c *fiber.Ctx) error {
	productId := c.Params("productId")
	// Get the product
	product := models.Product{}
	result := initializers.DB.First(&product, "id = ?", productId)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Product Not Found", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}

	result = initializers.DB.Delete(&models.Product{}, "id = ?", productId)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Product Not Deleted", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Product Deleted"})
}
