package api

import (
	"github.com/AzaanKH/hotel_reservation/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrUnAuthorizedD()
	}
	if !user.Admin {
		return ErrUnAuthorizedD()
	}
	return c.Next()
}
