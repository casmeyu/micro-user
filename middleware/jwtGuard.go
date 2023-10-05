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
		fmt.Println("Is a public route")
		return c.Next()
	}
	fmt.Println("Running JwtGuard")
	accessToken := c.GetReqHeaders()["Authorization"] // Authorization token must be `Authorization: Bearer <jwtToken>`
	accessToken = strings.Split(accessToken, " ")[1]  // Removes `Bearer ` part from Authorization token
	tokenData, err := auth.GetTokenData(accessToken)
	if err != nil {
		log.Println("[Middleware] (JwtGuard) - An error occurred in the JwtGuard", err)
		return c.Status(fiber.StatusUnauthorized).JSON("Unauthorized access")
	}
	fmt.Println("TOKEN DATA\n", tokenData)
	c.Locals("user", tokenData)
	fmt.Println(c.Locals("user"))
	return c.Next()
}
