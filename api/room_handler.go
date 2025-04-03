package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}


func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	// Validate that the booking dates are in the future and that FromDate is before TillDate.
	if err := params.validate(); err != nil {
		return err
	}
	
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON("error")
	}
	
	available, err := h.isRoomAvailableForBooking(c.Context(), roomID, params)
	if err != nil {
		return err
	}
	if !available {
		return fmt.Errorf("room already booked")
	}
	
	booking := types.Booking{
		RoomID:    roomID,
		UserID:    user.ID,
		FromDate:  params.FromDate,
		TillDate:  params.TillDate,
		NumPerson: params.NumPersons,
	}
	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	
	return c.JSON(inserted)
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	// Check that both booking dates are in the future.
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book a room in the past")
	}
	// Ensure the start date is before the end date.
	if p.FromDate.After(p.TillDate) || p.FromDate.Equal(p.TillDate) {
		return fmt.Errorf("fromDate must be before tillDate")
	}
	return nil
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error{
	rooms, err := h.store.Room.GetRooms(c.Context(),bson.M{})
	if err != nil{
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	// Find any booking that overlaps with the requested time range.
	filter := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$lt": params.TillDate, // existing booking starts before new booking ends
		},
		"tillDate": bson.M{
			"$gt": params.FromDate, // existing booking ends after new booking starts
		},
	}

	bookings, err := h.store.Booking.GetBookings(ctx, filter)
	if err != nil {
		return false, err
	}

	// If there are any overlapping bookings, the room is not available.
	return len(bookings) == 0, nil
}
