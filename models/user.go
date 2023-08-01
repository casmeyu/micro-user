package models

import (
	"log"
	"time"

	"github.com/casmeyu/micro-user/storage"
	"github.com/casmeyu/micro-user/structs"
	"github.com/casmeyu/micro-user/utils"
	"github.com/casmeyu/micro-user/validators"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string    `json:"username" gorm:"unique;not null;size:50" validate:"required,min=3,max=12"`
	Password       string    `json:"password" gorm:"not null" validate:"required_with_all"`
	RefreshToken   string    `json:"-"`
	LastConnection time.Time `json:"lastConnection"`
}

func HandleUserCreate(db *gorm.DB, c *fiber.Ctx) error {
	var err error
	var errors []*structs.IError
	var validate = validator.New()

	// Register custom validation for UserCreation
	validate.RegisterValidation("passwordRegex", validators.ValidatePasswordRegex)

	user := new(User)
	if err = c.BodyParser(user); err != nil {
		return c.Status(503).Send([]byte(err.Error()))
	}

	user.LastConnection = time.Now()
	err = validate.Struct(user)
	if err != nil {
		utils.FormatValidationErrors(err, &errors)
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	tx := db.Create(user)
	if tx.Error != nil {
		log.Println("[Users] (HandleUserCreate) - Error occurred while creating new user", tx.Error.Error())
		return c.Status(503).Send([]byte(tx.Error.Error()))
	}
	storage.Close(db)
	return nil
}
