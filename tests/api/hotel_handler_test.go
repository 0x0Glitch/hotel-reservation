package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0x0Glitch/hotel-reservation/api"
	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupHotelTest creates a test server with the hotel routes
func setupHotelTest(t *testing.T) (*fiber.App, *mongo.Client, func()) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Initialize stores and handlers
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)

	// Create a Store that wraps all stores
	store := &db.Store{
		Hotel: hotelStore,
		Room:  roomStore,
	}

	hotelHandler := api.NewHotelHandler(store)

	// Clean up collections before test
	collections := []string{"hotels", "rooms"}
	for _, coll := range collections {
		_, err = client.Database(db.DBNAME).Collection(coll).DeleteMany(context.TODO(), bson.M{})
		if err != nil {
			t.Fatalf("Error cleaning %s collection: %v", coll, err)
		}
	}

	// Setup Fiber app
	app := fiber.New()
	
	// Setup hotel routes manually
	app.Get("/api/hotels", hotelHandler.HandleGetHotels)
	app.Get("/api/hotels/:id", hotelHandler.HandleGetHotel)
	app.Get("/api/hotels/:id/rooms", hotelHandler.HandleGetRooms)

	// Return test server and cleanup function
	return app, client, func() {
		// Clean up after test
		for _, coll := range collections {
			_, err = client.Database(db.DBNAME).Collection(coll).DeleteMany(context.TODO(), bson.M{})
			if err != nil {
				t.Fatalf("Error cleaning %s collection: %v", coll, err)
			}
		}
		if err := client.Disconnect(context.TODO()); err != nil {
			t.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}
}

// insertTestHotels inserts test hotels into the database
func insertTestHotels(t *testing.T, client *mongo.Client) []*types.Hotel {
	hotelStore := db.NewMongoHotelStore(client)
	
	// Create test hotels
	hotels := []*types.Hotel{
		{
			ID:       primitive.NewObjectID(),
			Name:     "Luxury Hotel",
			Location: "Paris",
			Rooms:    []primitive.ObjectID{},
			Rating:   5,
		},
		{
			ID:       primitive.NewObjectID(),
			Name:     "Budget Inn",
			Location: "London",
			Rooms:    []primitive.ObjectID{},
			Rating:   3,
		},
	}
	
	// Insert hotels into database
	for _, hotel := range hotels {
		_, err := hotelStore.Insert(context.TODO(), hotel)
		if err != nil {
			t.Fatalf("Error inserting test hotel: %v", err)
		}
	}
	
	return hotels
}

// insertTestRooms inserts test rooms for a hotel
func insertTestRoom(t *testing.T, client *mongo.Client, hotelID primitive.ObjectID) *types.Room {
	roomStore := db.NewMongoRoomStore(client, db.NewMongoHotelStore(client))
	
	// Create test room
	room := &types.Room{
		Size:    "large",
		Seaside: true,
		Price:   199.99,
		HotelID: hotelID,
	}
	
	// Insert room
	insertedRoom, err := roomStore.InsertRoom(context.TODO(), room)
	if err != nil {
		t.Fatalf("Error inserting test room: %v", err)
	}
	
	return insertedRoom
}

// TestGetHotels tests fetching all hotels
func TestGetHotels(t *testing.T) {
	app, client, cleanup := setupHotelTest(t)
	defer cleanup()
	
	// Insert test hotels
	hotels := insertTestHotels(t, client)
	
	// Create HTTP request to get all hotels
	req := httptest.NewRequest(http.MethodGet, "/api/hotels", nil)
	
	// Execute request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
	
	// Parse response
	var fetchedHotels []types.Hotel
	if err := json.NewDecoder(resp.Body).Decode(&fetchedHotels); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	
	// Verify we got the expected number of hotels
	if len(fetchedHotels) != len(hotels) {
		t.Errorf("Expected %d hotels, got %d", len(hotels), len(fetchedHotels))
	}
}

// TestGetHotelByID tests fetching a specific hotel
func TestGetHotelByID(t *testing.T) {
	app, client, cleanup := setupHotelTest(t)
	defer cleanup()
	
	// Insert test hotels
	hotels := insertTestHotels(t, client)
	hotel := hotels[0] // Get first hotel
	
	// Create HTTP request to get the hotel
	req := httptest.NewRequest(http.MethodGet, "/api/hotels/"+hotel.ID.Hex(), nil)
	
	// Execute request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
	
	// Parse response
	var fetchedHotel types.Hotel
	if err := json.NewDecoder(resp.Body).Decode(&fetchedHotel); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	
	// Verify hotel fields
	if fetchedHotel.ID != hotel.ID {
		t.Errorf("Expected hotel ID %v, got %v", hotel.ID, fetchedHotel.ID)
	}
	if fetchedHotel.Name != hotel.Name {
		t.Errorf("Expected name %s, got %s", hotel.Name, fetchedHotel.Name)
	}
	if fetchedHotel.Location != hotel.Location {
		t.Errorf("Expected location %s, got %s", hotel.Location, fetchedHotel.Location)
	}
	if fetchedHotel.Rating != hotel.Rating {
		t.Errorf("Expected rating %d, got %d", hotel.Rating, fetchedHotel.Rating)
	}
}

// TestGetHotelRooms tests fetching rooms for a specific hotel
func TestGetHotelRooms(t *testing.T) {
	app, client, cleanup := setupHotelTest(t)
	defer cleanup()
	
	// Insert test hotels
	hotels := insertTestHotels(t, client)
	hotel := hotels[0] // Get first hotel
	
	// Insert a room for the hotel
	room := insertTestRoom(t, client, hotel.ID)
	
	// Create HTTP request to get hotel rooms
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/hotels/%s/rooms", hotel.ID.Hex()), nil)
	
	// Execute request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
	
	// Parse response
	var fetchedRooms []types.Room
	if err := json.NewDecoder(resp.Body).Decode(&fetchedRooms); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	
	// Verify we got at least one room
	if len(fetchedRooms) == 0 {
		t.Errorf("Expected at least one room, got none")
	}
	
	// Verify room fields
	if len(fetchedRooms) > 0 {
		fetchedRoom := fetchedRooms[0]
		if fetchedRoom.ID != room.ID {
			t.Errorf("Expected room ID %v, got %v", room.ID, fetchedRoom.ID)
		}
		if fetchedRoom.Size != room.Size {
			t.Errorf("Expected size %s, got %s", room.Size, fetchedRoom.Size)
		}
		if fetchedRoom.Seaside != room.Seaside {
			t.Errorf("Expected seaside %v, got %v", room.Seaside, fetchedRoom.Seaside)
		}
		if fetchedRoom.Price != room.Price {
			t.Errorf("Expected price %.2f, got %.2f", room.Price, fetchedRoom.Price)
		}
		if fetchedRoom.HotelID != hotel.ID {
			t.Errorf("Expected hotel ID %v, got %v", hotel.ID, fetchedRoom.HotelID)
		}
	}
} 