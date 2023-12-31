package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/casmeyu/micro-user/auth"
	"github.com/casmeyu/micro-user/configuration"
	"github.com/casmeyu/micro-user/middleware"
	"github.com/casmeyu/micro-user/storage"
	"github.com/casmeyu/micro-user/structs"
	"github.com/casmeyu/micro-user/userService"
	"github.com/casmeyu/micro-user/utils"
	"github.com/casmeyu/micro-user/validators"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Setting config as global variable
var Config structs.Config
var Validate = validator.New()
var Db *gorm.DB

// Executes LoadConfig() function and sets up initial information for the backend app
func Setup() error {
	err := configuration.LoadConfig(&Config)
	if err != nil {
		log.Println("[Main] (Setup) - Error while setting up config", err.Error())
		return err
	}
	return nil
}

func SetRoutes(app *fiber.App) {
	// Setup User Routes
	userRoutes := app.Group("/users")
	userRoutes.Get("/", func(c *fiber.Ctx) error {
		var dbUsers []userService.User
		var resUsers []structs.PublicUser
		if tx := Db.Find(&dbUsers); tx.Error != nil {
			log.Println("[GET] (/users) - Error occurred while finding users", tx.Error.Error())
			c.Status(501).JSON("Error while getting users")
		}
		for _, user := range dbUsers {
			resUsers = append(resUsers, structs.PublicUser{
				Id:             user.ID,
				Username:       user.Username,
				LastConnection: user.LastConnection,
			})
		}

		return c.Status(200).JSON(resUsers)
	})

	userRoutes.Post("/", func(c *fiber.Ctx) error {
		var err error
		var errors []*structs.IError

		user := new(userService.User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(503).Send([]byte(err.Error()))
		}
		err = Validate.Struct(user)
		if err != nil {
			utils.FormatValidationErrors(err, &errors)
			return c.Status(fiber.StatusBadRequest).JSON(errors)
		}
		Db, err := storage.Open(Config) // Pass only Config.Db as it is more clean and efficient
		if err != nil {
			log.Println("[POST] (/users) - Error trying to connect to database", err.Error())
		}

		res := userService.CreateUser(user, Db)

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

		res := userService.GetById(userId, Db)

		if res.Success == true {
			return c.Status(res.Status).JSON(res.Result)
		} else {
			return c.Status(res.Status).JSON(res.Err)
		}
	})
	// END User Routes

	app.Use("/login", middleware.IsPublic) // Setting login as a public route
	app.Post("/login", func(c *fiber.Ctx) error {
		var errors []*structs.IError
		var err error
		userLogin := new(structs.UserLogin)
		if err := c.BodyParser(userLogin); err != nil {
			log.Println("[Auth] (Login) Error occurred while parsing request body", err.Error())
			c.Status(503).SendString("Error while parsing body request")
		}

		err = Validate.Struct(userLogin)
		if err != nil {
			utils.FormatValidationErrors(err, &errors)
			return c.Status(fiber.StatusBadRequest).JSON(errors)
		}
		Db, err := storage.Open(Config) // Pass only Config.Db as it is more clean and efficient
		if err != nil {
			log.Println("[POST] (/login) - Error trying to connect to database", err.Error())
		}

		res := auth.Login(userLogin, Db)

		if res.Success {
			return c.Status(res.Status).JSON(res.Result)
		} else {
			return c.Status(res.Status).JSON(res.Err)
		}
	})

	app.Post("/logout", func(c *fiber.Ctx) error {
		// Get JwtToken
		// Check JwtToken agains Db.users
		// Somehow invalidate JwtToken and RefreshToken
		return nil
	})

	// Test for "private" route with JwtGuard middleware
	app.Get("/private", func(c *fiber.Ctx) error {
		fmt.Println("Running private route")
		fmt.Println(c.Locals("user"))
		return nil
	})

	// Add App middleware to be run after specific route middleware
	app.Use(middleware.JwtGuard)
}

func main() {
	Validate.RegisterValidation("passwordRegex", validators.ValidatePasswordRegex)
	if err := Setup(); err != nil {
		os.Exit(2)
	}
	log.Println("Configuration loaded")

	storage.MakeMigration(Config, &userService.User{})

	app := fiber.New()
	SetRoutes(app)

	conn, err := storage.Open(Config) // Pass only Config.Db as it is more clean and efficient
	if err != nil {
		os.Exit(2)
	}
	Db = conn // Setting DB Connection for all the routes
	log.Printf("Connected to %s database: %s\n", Db.Name(), Config.Db.Name)
	app.Listen(":3000")
}
