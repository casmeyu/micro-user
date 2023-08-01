package auth

import (
	"fmt"
	"log"

	"github.com/casmeyu/micro-user/models"
	"github.com/casmeyu/micro-user/structs"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(db *gorm.DB, c *fiber.Ctx) error {
	var err error
	userLogin := new(structs.UserLogin)
	if err := c.BodyParser(userLogin); err != nil {
		log.Println("[Auth] (Login) Error occurred while parsing request body", err.Error())
		return c.Status(503).Send([]byte(err.Error()))
	}

	var dbUser = models.User{Username: userLogin.Username}

	tx := db.First(&dbUser)
	if tx.Error != nil {
		log.Println("[Auth] (Login) - Error occurred while trying to get user", tx.Error.Error())
		c.Status(501).Send([]byte(tx.Error.Error()))
	}
	// TODO: create util functions: hashUserPwd and compareUserPwd
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(userLogin.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[Auth] (Login) - Error occurred while hashing userLogin password", err.Error())
	}

	err = bcrypt.CompareHashAndPassword(hashPwd, []byte(userLogin.Password))
	if err != nil {
		fmt.Println(err)
	} else {
		// WARNING: user password is being send in the response, remove it
		c.Status(fiber.StatusOK).JSON(dbUser)
	}
	return nil
}
