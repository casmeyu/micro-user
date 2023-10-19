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
		log.Println("Error while setting up config", err.Error())
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
		storage.Close(Db)
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
		storage.Close(Db)
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
		storage.Close(Db)
		if res.Success == true {
			return c.Status(res.Status).JSON(res.Result)
		} else {
			return c.Status(res.Status).JSON(res.Err)
		}
	})
	// END User Routes

	// LogIn and Logout
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
		storage.Close(Db)
		if res.Success == true {
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
	// END Login/Logout

	app.Use("/private", middleware.JwtGuard)
	// Test for "private" route (jwt middleware)
	app.Get("/private", func(c *fiber.Ctx) error {
		fmt.Println("Running private route")
		fmt.Println(c.Locals("user"))
		return nil
	})
}

func main() {
	Validate.RegisterValidation("passwordRegex", validators.ValidatePasswordRegex)
	if err := Setup(); err != nil {
		os.Exit(2)
	}
	log.Println("Configuration loaded")

	// Should I use a SINGLE Database Connection for the service?
	// Or each Route should connect and disconnect from the DB each time they are called???
	conn, err := storage.Open(Config) // Pass only Config.Db as it is more clean and efficient
	if err != nil {
		os.Exit(2)
	}
	Db = conn // Setting DB Connection for all the routes
	log.Printf("Connected to %s database: %s\n", Db.Name(), Config.Db.Name)
	storage.MakeMigration(Config, &userService.User{})

	app := fiber.New()

	SetRoutes(app)

	app.Listen(":3000")
}
