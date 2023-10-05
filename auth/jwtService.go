package auth

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	exp, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	if err != nil {
		log.Println("[Auth] (CreateJwtToken) - Error occurred while getting Jwt Expiration ENV", err.Error())
		tokenClaims["exp"] = time.Now().Add(time.Hour).Unix()
	} else {
		tokenClaims["exp"] = time.Now().Add(time.Second * time.Duration(exp)).Unix()
	}

	s, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println("[Auth] (CreateJwtToken) - Error occurred while signing JWT token", err.Error())
	}
	return s, err
}

func GetTokenData(tokenString string) (map[string]interface{}, error) {
	fmt.Println("Get Token data from", tokenString)
	var data map[string]interface{}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	fmt.Println(token)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims["data"] != nil {
		data = claims["data"].(map[string]interface{})
	} else {
		data = nil
	}

	return data, nil
}
