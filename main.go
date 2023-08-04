package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/casmeyu/micro-user/auth"
	"github.com/casmeyu/micro-user/configuration"
	"github.com/casmeyu/micro-user/storage"
	"github.com/casmeyu/micro-user/structs"
	"github.com/casmeyu/micro-user/userService"
	"github.com/casmeyu/micro-user/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Setting config as global variable
var Config structs.Config
var Validator = validator.New()

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
	// Setup User Routes
	userRoutes := app.Group("/users")

	userRoutes.Get("/", func(c *fiber.Ctx) error {
		db, err := storage.Open(Config)
		if err != nil {
			log.Println("[GET] (/users) - Error trying to connect to database", err.Error())
		}
		var users []userService.User
		db.Find(&users)
		storage.Close(db)
		return c.JSON(users)
	})

	userRoutes.Post("/", func(c *fiber.Ctx) error {
		user := new(userService.User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(503).Send([]byte(err.Error()))
		}
		db, err := storage.Open(Config)
		if err != nil {
			log.Println("[POST] (/users) - Error trying to connect to database", err.Error())
		}
		res := userService.CreateUser(user, db)
		if res.Success == true {
			return c.Status(res.Status).JSON(res.Result)
		} else {
			return c.Status(res.Status).JSON(res.Err)
		}
	})

	userRoutes.Get("/:id<int>", func(c *fiber.Ctx) error {
		var err error
		userId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Invalid user id")
		}
		db, err := storage.Open(Config)
		if err != nil {
			log.Println("[POST] (/users) - Error trying to connect to database", err.Error())
		}
		res := userService.GetById(userId, db)
		if res.Success == true {
			return c.Status(res.Status).JSON(res.Result)
		} else {
			return c.Status(res.Status).JSON(res.Err)
		}
	})
	// END User Routes

	app.Post("/login", func(c *fiber.Ctx) error {
		var errors []*structs.IError
		var err error
		userLogin := new(structs.UserLogin)
		if err := c.BodyParser(userLogin); err != nil {
			log.Println("[Auth] (Login) Error occurred while parsing request body", err.Error())
			c.Status(503).SendString("Error while parsing body request")
		}
		// Validate userLogin request
		err = Validator.Struct(userLogin)
		if err != nil {
			utils.FormatValidationErrors(err, &errors)
			return c.Status(fiber.StatusBadRequest).JSON(errors)
		}
		db, err := storage.Open(Config)
		if err != nil {
			log.Println("[POST] (/login) - Error trying to connect to database", err.Error())
		}

		res := auth.Login(userLogin, db)
		if res.Success == true {
			return c.Status(res.Status).JSON(res.Result)
		} else {
			return c.Status(res.Status).JSON(res.Err)
		}
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
	storage.MakeMigration(Config, &userService.User{})

	app := fiber.New()

	SetRoutes(app)

	app.Listen(":3000")
}
