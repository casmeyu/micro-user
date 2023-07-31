package main

import (
	"fmt"
	"log"
	"os"

	"github.com/casmeyu/micro-user/configuration"
	"github.com/casmeyu/micro-user/models"
	"github.com/casmeyu/micro-user/storage"
	"github.com/casmeyu/micro-user/structs"
	"github.com/gofiber/fiber/v2"
)

// Setting config as global variable
var Config structs.Config

// Executes LoadConfig() function and sets up initial information for the backend app
func Setup() error {
	err := configuration.LoadConfig(&Config)
	if err != nil {
		log.Println("Error while setting up config", err.Error())
		return err
	}
	return nil
}

func SetRoutes(app *fiber.App) {
	app.Get("/users", func(c *fiber.Ctx) error {
		db, err := storage.Open(Config)
		if err != nil {
			log.Println("[GET] (/users) - Error occurred while getting users", err.Error())
		}
		var users []models.User
		db.Find(&users)
		storage.Close(db)
		return c.JSON(users)
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		db, err := storage.Open(Config)
		if err != nil {
			log.Println("[POST] (/users) - Error occurres while creating a user", err.Error())
		}
		storage.Close(db)
		return models.HandleUserCreate(db, c)
	})
}

func main() {
	fiber.New()
	if err := Setup(); err != nil {
		os.Exit(2)
	}
	log.Println("Configuration loaded")

	db, err := storage.Open(Config)
	if err != nil {
		os.Exit(2)
	}
	fmt.Println("Db is connected!", &db)
	storage.MakeMigration(Config, &models.User{})

	app := fiber.New()

	SetRoutes(app)

	app.Listen(":3000")
}
