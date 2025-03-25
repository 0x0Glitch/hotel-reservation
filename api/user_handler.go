package api

import (
	"github.com/0x0Glitch/hotel-reservation/types"

	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(c *fiber.Ctx) error{
	u := types.User{
		FirstName: "Anshuman",
		LastName: "Yadav",
	}

	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error{
	return c.JSON("HandleUser->Anshuman")
}