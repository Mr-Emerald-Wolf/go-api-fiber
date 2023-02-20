package main

import (
	"github.com/mr-emerald-wolf/go-api-fiber/initializers"
	"github.com/mr-emerald-wolf/go-api-fiber/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Product{})

}
