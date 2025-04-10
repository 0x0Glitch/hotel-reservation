package api

import (
	"net/http"

	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct{
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler{
	return &BookingHandler{
		store:store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error{
	bookings, err := h.store.Booking.GetBookings(c.Context(),bson.M{})
	if err != nil{
		return err
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error{
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(),id)
	if err != nil{
		return err
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok{
		return err
	}
	if booking.UserID != user.ID{
	return c.Status(http.StatusUnauthorized).JSON("errrrorororororororo")
}
	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	
	// First, get the booking to check if user is authorized
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "booking not found",
		})
	}
	
	// Get the user from context (set by JWT middleware)
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	
	// Check if user is authorized (either the booking owner or an admin)
	if booking.UserID != user.ID && !user.IsAdmin {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "you do not have permission to cancel this booking",
		})
	}
	
	// Cancel the booking
	if err := h.store.Booking.CancelBooking(c.Context(), id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to cancel booking",
		})
	}
	
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "booking successfully cancelled",
	})
}


