package middleware

import (
	"log"
	"strings"

	"github.com/casmeyu/micro-user/auth"
	"github.com/casmeyu/micro-user/constants"
	"github.com/gofiber/fiber/v2"
)

func JwtGuard(c *fiber.Ctx) error {
	if c.Locals(constants.IS_PUBLIC) == true {
		return c.Next()
	}
	accessToken := c.GetReqHeaders()["Authorization"] // Authorization token must be `Authorization: Bearer <jwtToken>`
	accessToken = strings.Split(accessToken, " ")[1]  // Removes `Bearer ` part from Authorization token
	tokenData, err := auth.GetTokenData(accessToken)
	if err != nil {
		log.Println("[Middleware] (JwtGuard) - An error occurred in the JwtGuard", err)
		return c.Status(fiber.StatusUnauthorized).JSON("Unauthorized access")
	}

	c.Locals("user", tokenData)
	return c.Next()
}
