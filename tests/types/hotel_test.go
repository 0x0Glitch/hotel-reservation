package types

import (
	"testing"

	"github.com/0x0Glitch/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestHotelStructure validates that the Hotel struct can be properly instantiated
// and all fields are correctly set and retrieved
func TestHotelStructure(t *testing.T) {
	// Set up test data
	hotelID := primitive.NewObjectID()
	roomIDs := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}
	
	// Create a new hotel
	hotel := types.Hotel{
		ID:       hotelID,
		Name:     "Test Hotel",
		Location: "Test Location",
		Rooms:    roomIDs,
		Rating:   5,
	}
	
	// Validate ID field
	if hotel.ID != hotelID {
		t.Errorf("expected ID %v, got %v", hotelID, hotel.ID)
	}
	
	// Validate Name field
	if hotel.Name != "Test Hotel" {
		t.Errorf("expected Name 'Test Hotel', got '%s'", hotel.Name)
	}
	
	// Validate Location field
	if hotel.Location != "Test Location" {
		t.Errorf("expected Location 'Test Location', got '%s'", hotel.Location)
	}
	
	// Check rooms array length
	if len(hotel.Rooms) != len(roomIDs) {
		t.Errorf("expected %d rooms, got %d", len(roomIDs), len(hotel.Rooms))
	}
	
	// Validate all room IDs match
	for i, roomID := range roomIDs {
		if hotel.Rooms[i] != roomID {
			t.Errorf("Room ID mismatch at index %d: expected %v, got %v", i, roomID, hotel.Rooms[i])
		}
	}
	
	// Validate Rating field
	if hotel.Rating != 5 {
		t.Errorf("expected Rating 5, got %d", hotel.Rating)
	}
}

// TestRoomStructure validates that the Room struct can be properly instantiated
// and all fields are correctly set and retrieved
func TestRoomStructure(t *testing.T) {
	// Set up test data
	roomID := primitive.NewObjectID()
	hotelID := primitive.NewObjectID()
	
	// Create a new room
	room := types.Room{
		ID:      roomID,
		Seaside: true,
		Size:    "Double",
		Price:   150.50,
		HotelID: hotelID,
	}
	
	// Validate ID field
	if room.ID != roomID {
		t.Errorf("expected ID %v, got %v", roomID, room.ID)
	}
	
	// Validate Seaside field
	if !room.Seaside {
		t.Errorf("expected Seaside to be true")
	}
	
	// Validate Size field
	if room.Size != "Double" {
		t.Errorf("expected Size 'Double', got '%s'", room.Size)
	}
	
	// Validate Price field
	if room.Price != 150.50 {
		t.Errorf("expected Price 150.50, got %f", room.Price)
	}
	
	// Validate HotelID field
	if room.HotelID != hotelID {
		t.Errorf("expected HotelID %v, got %v", hotelID, room.HotelID)
	}
}

// TestRoomTypes validates the room type constants
func TestRoomTypes(t *testing.T) {
	// Validate room type constant values
	if types.SingleRoomType != 1 {
		t.Errorf("expected SingleRoomType to be 1, got %d", types.SingleRoomType)
	}
	
	if types.DoubleRoomType != 2 {
		t.Errorf("expected DoubleRoomType to be 2, got %d", types.DoubleRoomType)
	}
	
	if types.SeaSideRoomType != 3 {
		t.Errorf("expected SeaSideRoomType to be 3, got %d", types.SeaSideRoomType)
	}
	
	if types.DeluxRoomType != 4 {
		t.Errorf("expected DeluxRoomType to be 4, got %d", types.DeluxRoomType)
	}
} 