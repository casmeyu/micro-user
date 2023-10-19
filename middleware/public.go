package middleware

import (
	"github.com/casmeyu/micro-user/constants"
	"github.com/gofiber/fiber/v2"
)

func IsPublic(c *fiber.Ctx) error {
	c.Locals(constants.IS_PUBLIC, true)
	return c.Next()
}
