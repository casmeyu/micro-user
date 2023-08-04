package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Custom validation for user password with a regex
// RIGHT NOW IT DOES NOT GET CALLED
func ValidatePasswordRegex(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	regex := `^[0-9]$` // THIS SHOULD BE IN A ENVIRONMENT VARIABLE

	match, _ := regexp.MatchString(regex, password)
	return match
}
