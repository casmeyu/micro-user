package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJwtToken(claims map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)

	for key := range claims {
		tokenClaims[key] = claims[key]
	}
	// Add expiration to JWT
	// exp, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	// if err != nil {
	// 	log.Println("[Auth] (CreateJwtToken) - Error occurred while getting Jwt Expiration ENV", err.Error())
	// 	tokenClaims["exp"] = time.Hour
	// } else {
	tokenClaims["exp"] = time.Hour
	// }

	s, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println("[Auth] (CreateJwtToken) - Error occurred while signing JWT token", err.Error())
	}
	return s, err
}
