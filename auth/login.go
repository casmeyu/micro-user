package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/casmeyu/micro-user/structs"
	"github.com/casmeyu/micro-user/userService"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var Validator = validator.New()

func Login(userLogin *structs.UserLogin, db *gorm.DB) structs.ServiceResponse {
	var err error
	var tx *gorm.DB
	var res = structs.ServiceResponse{}
	var dbUser userService.User
	var resUser structs.PublicUser

	tx = db.Model(&userService.User{}).Where("username=?", userLogin.Username).First(&dbUser)
	if tx.Error != nil {
		log.Println("[Auth] (Login) - Error occurred while trying to get user:", tx.Error.Error())
		res.Err = fmt.Sprintf("User: %s", tx.Error.Error())
		res.Status = fiber.StatusInternalServerError
		return res
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(userLogin.Password))
	if err != nil {
		res.Err = "Wrong username or password"
		res.Status = fiber.StatusBadRequest
		return res
	} else {
		dbUser.LastConnection = time.Now()
		tx = db.Save(&dbUser)
		if tx.Error != nil {
			log.Println("[Users] (CreateUser) - Error occurred while parsing the user to Json", tx.Error.Error())
			res.Err = "Error while login in"
			res.Status = fiber.StatusInternalServerError
			return res
		}

		resUser = structs.PublicUser{
			Id:             dbUser.ID,
			Username:       dbUser.Username,
			LastConnection: dbUser.LastConnection,
		}
		res.Success = true
		res.Result = resUser
		res.Status = fiber.StatusOK
		return res
	}
}
