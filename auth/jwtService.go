package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJwtToken(claims map[string]interface{}) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
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

func GetTokenData(tokenString string) (map[string]interface{}, error) {
	fmt.Println("GET TOKEN DATA")
	var data map[string]interface{}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	claims := token.Claims.(jwt.MapClaims)
	if claims["data"] != nil {
		data = claims["data"].(map[string]interface{})
	} else {
		data = map[string]interface{}{}
	}

	return data, nil
}
