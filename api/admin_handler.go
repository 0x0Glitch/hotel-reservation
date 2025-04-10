package api

import (
	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminHandler struct {
	store *db.Store
}

func NewAdminHandler(store *db.Store) *AdminHandler {
	return &AdminHandler{
		store: store,
	}
}

func (h *AdminHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.store.User.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *AdminHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (h *AdminHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *AdminHandler) HandleMakeAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	
	// Check if ID is a valid ObjectID format
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID format",
		})
	}
	
	// Get the user
	user, err := h.store.User.GetUserById(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}
	
	// Update user to make them admin
	user.IsAdmin = true
	
	// Update the user in the database
	if err := h.store.User.UpdateUser(c.Context(), bson.M{"_id": user.ID}, bson.M{"isAdmin": true}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update user",
		})
	}
	
	return c.JSON(user)
} 