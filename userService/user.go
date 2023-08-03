package userService

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string    `json:"username" gorm:"unique;not null;size:50"`
	Password       string    `json:"password" gorm:"not null"`
	RefreshToken   string    `json:"-"`
	LastConnection time.Time `json:"lastConnection"`
}

func HandleUserCreate(db *gorm.DB, c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(503).Send([]byte(err.Error()))
	}

	user.LastConnection = time.Now()
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[Users] (HandleUserCreate) - Error occurred while hashing user's password")
	}
	user.Password = string(hashPwd)
	// TODO: Add entity validations and set default last_connection to null in struct

	tx := db.Create(user)
	if tx.Error != nil {
		log.Println("[Users] (HandleUserCreate) - Error occurred while writing user to the database", tx.Error.Error())
		return c.Status(503).Send([]byte(tx.Error.Error()))
	}

	return c.Status(201).JSON(user)
}
