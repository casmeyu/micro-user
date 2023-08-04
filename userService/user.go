package userService

import (
	"fmt"
	"log"
	"time"

	"github.com/casmeyu/micro-user/structs"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string    `json:"username" gorm:"unique;not null;size:50"`
	Password       string    `json:"password" gorm:"not null" validate:"passwordRegex"`
	RefreshToken   string    `json:"refreshToken"`
	LastConnection time.Time `json:"lastConnection"`
}

func CreateUser(user *User, db *gorm.DB) structs.ServiceResponse {
	var res = structs.ServiceResponse{}

	user.LastConnection = time.Now()
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[Users] (CreateUser) - Error occurred while hashing user's password")
	}
	user.Password = string(hashPwd)
	// TODO: Add entity validations and set default last_connection to null in struct

	tx := db.Create(user)
	if tx.Error != nil {
		log.Println("[Users] (CreateUser) - Error occurred while writing user to the database", tx.Error.Error())
		var dupErrMsg = fmt.Sprintf("Error 1062 (23000): Duplicate entry '%s' for key 'users.username'", user.Username)
		if tx.Error.Error() == dupErrMsg {
			res.Err = fmt.Sprintf("Username %s is already taken", user.Username)
			res.Status = fiber.StatusBadRequest
		} else {
			res.Err = "Error occurred while creating new user"
			res.Status = 503
		}
		return res
	}
	res.Success = true
	res.Result = user
	res.Status = fiber.StatusCreated
	return res
}

func GetById(id int, db *gorm.DB) structs.ServiceResponse {
	res := structs.ServiceResponse{}
	var dbUser User
	tx := db.Model(&User{}).Where("id=?", id).First(&dbUser)
	if tx.Error != nil {
		log.Println("[Auth] (Login) - Error occurred while trying to get user:", tx.Error.Error())
		res.Err = "User does not exist"
		res.Status = fiber.StatusNotFound
		return res
	}
	res.Success = true
	res.Status = fiber.StatusOK
	res.Result = dbUser
	return res
}
