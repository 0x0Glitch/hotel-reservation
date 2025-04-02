package api

import (
	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HotelHandler handles HTTP requests related to hotel and room operations
// It processes requests for listing hotels, getting hotel details, and viewing rooms
type HotelHandler struct{
	store 	*db.Store  // Central store providing access to all database collections
}

// NewHotelHandler creates a new HotelHandler with the provided store
// Factory function to create handlers with dependency injection
func NewHotelHandler(store *db.Store)*HotelHandler{
	return &HotelHandler{
		store: store,
	}
}

// type HotelQueryParams struct{
// 	Rooms bool
// 	Rating int
// }

// HandleGetRooms processes requests to get all rooms for a specific hotel
// GET /api/hotel/:id/rooms
func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error{
	// Extract hotel ID from URL parameters
	id := c.Params("id")
	
	// Convert string ID to MongoDB ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		return err
	}
	
	// Create filter to find rooms for this specific hotel
	filter := bson.M{"hotelID": oid}
	
	// Fetch rooms from the database
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil{
		return err
	}
	
	// Return rooms as JSON array
	return c.JSON(rooms)
}

// HandleGetHotels processes requests to get all hotels
// GET /api/hotel
// Can be extended to support query parameters for filtering
func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error{
	// The nil filter means "get all hotels" (no conditions)
	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil)
	if err != nil{
		return err
	}
	
	// Return hotels as JSON array
	return c.JSON(hotels)
}

// HandleGetHotel processes requests to get a specific hotel by ID
// GET /api/hotel/:id
func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error{
	// Extract hotel ID from URL parameters
	id := c.Params("id")
	
	// Convert string ID to MongoDB ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		return err
	}
	
	// Fetch the hotel from the database
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), oid)
	if err != nil {
		return err
	}
	
	// Return the hotel as JSON
	return c.JSON(hotel)
}
