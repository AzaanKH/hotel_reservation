package middleware

import (
	"fmt"

	"github.com/AzaanKH/hotel_reservation/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("not authorized")
	}
	if !user.Admin {
		return fmt.Errorf("not authorized")
	}
	return c.Next()
}
