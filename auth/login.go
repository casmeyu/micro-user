package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/casmeyu/micro-user/structs"
	"github.com/casmeyu/micro-user/userService"
	"github.com/casmeyu/micro-user/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var Validator = validator.New()

func Login(db *gorm.DB, c *fiber.Ctx) error {
	var err error
	var errors []*structs.IError
	var tx *gorm.DB
	var res = structs.ServiceResponse{
		Success: false,
		Result:  nil,
		Err:     "",
	}
	var dbUser userService.User

	userLogin := new(structs.UserLogin)
	if err := c.BodyParser(userLogin); err != nil {
		res.Err = fmt.Sprintf("User: %s", err.Error())
		log.Println("[Auth] (Login) Error occurred while parsing request body", err.Error())
		return c.Status(503).JSON(res)
	}

	// UserLogin Input validation
	err = Validator.Struct(userLogin)
	if err != nil {
		utils.FormatValidationErrors(err, &errors)
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	tx = db.Model(&userService.User{}).Where("username=?", userLogin.Username).First(&dbUser)
	if tx.Error != nil {
		log.Println("[Auth] (Login) - Error occurred while trying to get user:", tx.Error.Error())
		res.Err = fmt.Sprintf("User: %s", tx.Error.Error())
		return c.Status(501).JSON(res)
	}
	// TODO: create util functions: hashUserPwd and compareUserPwd
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(userLogin.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[Auth] (Login) - Error occurred while hashing userLogin password", err.Error())
	}

	err = bcrypt.CompareHashAndPassword(hashPwd, []byte(userLogin.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Wrong credentials")
	} else {
		dbUser.LastConnection = time.Now()
		tx = db.Save(&dbUser)
		if tx.Error != nil {
			log.Println("[Users] (HandleUserCreate) - Error occurred while parsing the user to Json", tx.Error.Error())
			return c.Status(501).SendString("Error while loging in")
		}
		// WARNING: user password is being send in the response, remove it
		return c.Status(fiber.StatusOK).JSON(dbUser)
	}
}
