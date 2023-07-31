package models

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
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
	// TODO: Add entity validations and set default last_connection to null in struct

	tx := db.Create(user)
	if tx.Error != nil {
		log.Println("[Users] (HandleUserCreate) - Error occurred while creating new user", tx.Error.Error())
		return c.Status(503).Send([]byte(tx.Error.Error()))
	}

	return c.Status(201).JSON(user)
}
