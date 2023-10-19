package middleware

import (
	"fmt"
	"log"
	"strings"

	"github.com/casmeyu/micro-user/auth"
	"github.com/gofiber/fiber/v2"
)

func JwtGuard(c *fiber.Ctx) error {
	if c.Locals("IS_PUBLIC") == true {
		return c.Next()
	}
	fmt.Println("Running JwtGuard")
	accessToken := c.Get("Authorization") // Authorization token must be `Authorization: Bearer <jwtToken>`
	if accessToken == "" {
		log.Println("[Middleware] (JwtGuard) - An error occurred in the JwtGuard", "No Authorization token was provided")
		return c.Status(fiber.StatusUnauthorized).JSON("Unauthorized access")
	}
	accessToken = strings.Split(accessToken, " ")[1] // Removes `Bearer ` part from Authorization token
	tokenData, err := auth.GetTokenData(accessToken)
	if err != nil {
		log.Println("[Middleware] (JwtGuard) - An error occurred in the JwtGuard", err)
		return c.Status(fiber.StatusUnauthorized).JSON("Unauthorized access")
	}
	c.Locals("user", tokenData)
	return c.Next()
}
