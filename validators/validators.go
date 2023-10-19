package validators

import (
	"log"
	"os"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Custom validation for user password with a regex
func ValidatePasswordRegex(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	regex, exist := os.LookupEnv("USER_PASSWORD_REGEX")
	if !exist {
		regex = "^.{3,}$"
	} // THIS SHOULD BE IN A ENVIRONMENT VARIABLE

	match, err := regexp.MatchString(regex, password)
	if err != nil {
		log.Println("[Validator] (ValidatePasswordRegex) - Error occurred matching regex", err.Error())
	}
	return match
}
