package api

import (
	"context"

	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"

	"github.com/gofiber/fiber/v2"
)
type UserHandler struct{
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler{
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error{
	var
	( id = c.Params("id")
	ctx = context.Background()
)
	user, err:= h.userStore.GetUserById(ctx,id)
	if err != nil{
		return err
	}
	return c.JSON(user)
	return c.JSON("HandleUser->Anshuman")
}


func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error{
	u := types.User{
		FirstName: "Anshuman",
		LastName: "Yadav",
	}

	return c.JSON(u)
}