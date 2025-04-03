package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct{
	FromDate time.Time  `json:"fromDate"`
	TillDate time.Time  `json:"tillDate"`
	NumPersons int      `json"numPersons`
}

type RoomHandler struct{
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler{

	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error{
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil{
		return err
	}
	roomID,err := primitive.ObjectIDFromHex(c.Params("id"))
	
	if err != nil{
		return err
	}
	// return nil
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON("error")
	}
	booking := types.Booking{
		RoomID: roomID,
		UserID:user.ID,
		FromDate: params.FromDate,
		TillDate: params.TillDate,
		NumPerson: params.NumPersons,
	}
	fmt.Printf("%+v",booking)
	return nil


}

