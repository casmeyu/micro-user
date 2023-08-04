package validators

import (
	"log"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Custom validation for user password with a regex
// RIGHT NOW IT DOES NOT GET CALLED
func ValidatePasswordRegex(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	regex := `^.{3,}$` // THIS SHOULD BE IN A ENVIRONMENT VARIABLE

	match, err := regexp.MatchString(regex, password)
	if err != nil {
		log.Println("Error occurred matching regex", err.Error())
	}
	return match
}
