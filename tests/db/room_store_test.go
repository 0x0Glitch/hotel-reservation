package db

import (
	"context"
	"testing"

	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testDBURI  = "mongodb://localhost:27017"
)

// roomTestDB wraps the room store and client for testing
type roomTestDB struct {
	client    *mongo.Client
	roomStore db.RoomStore
	hotelStore db.HotelStore
}

// setupRoomTest initializes the test environment with a MongoDB connection
func setupRoomTest(t *testing.T) *roomTestDB {
	// Connect to test database
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDBURI))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	
	// Clean up any previous test data to ensure a fresh start
	_, err = client.Database(db.DBNAME).Collection("rooms").DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		t.Fatalf("Error cleaning up rooms collection: %v", err)
	}
	
	_, err = client.Database(db.DBNAME).Collection("hotels").DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		t.Fatalf("Error cleaning up hotels collection: %v", err)
	}

	return &roomTestDB{
		client:    client,
		roomStore: roomStore,
		hotelStore: hotelStore,
	}
}

// teardown cleans up after tests
func (tdb *roomTestDB) teardown(t *testing.T) {
	// Delete all test data
	_, err := tdb.client.Database(db.DBNAME).Collection("rooms").DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		t.Fatalf("Error cleaning up rooms collection: %v", err)
	}
	
	_, err = tdb.client.Database(db.DBNAME).Collection("hotels").DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		t.Fatalf("Error cleaning up hotels collection: %v", err)
	}

	// Disconnect from MongoDB
	if err := tdb.client.Disconnect(context.TODO()); err != nil {
		t.Fatalf("Error disconnecting from MongoDB: %v", err)
	}
}

// TestMongoRoomStore_InsertRoom tests inserting a new room
func TestMongoRoomStore_InsertRoom(t *testing.T) {
	// Skip integration tests when running in short mode
	if testing.Short() {
		t.Skip("Skipping MongoDB integration test in short mode")
	}

	tdb := setupRoomTest(t)
	defer tdb.teardown(t)

	// First create a hotel to associate the room with
	hotel := &types.Hotel{
		Name:     "Test Hotel",
		Location: "Test Location",
		Rooms:    []primitive.ObjectID{},
		Rating:   4,
	}

	insertedHotel, err := tdb.hotelStore.Insert(context.TODO(), hotel)
	if err != nil {
		t.Fatalf("error inserting hotel: %v", err)
	}

	// Create a test room
	room := &types.Room{
		Size:    "large",
		Seaside: true,
		Price:   129.99,
		HotelID: insertedHotel.ID,
	}

	// Insert the room
	insertedRoom, err := tdb.roomStore.InsertRoom(context.TODO(), room)
	if err != nil {
		t.Fatalf("error inserting room: %v", err)
	}

	// Verify room ID is set
	if insertedRoom.ID.IsZero() {
		t.Errorf("expected room ID to be set")
	}

	// Verify room fields match
	if insertedRoom.Size != room.Size {
		t.Errorf("expected Size %s, got %s", room.Size, insertedRoom.Size)
	}

	if insertedRoom.Seaside != room.Seaside {
		t.Errorf("expected Seaside %v, got %v", room.Seaside, insertedRoom.Seaside)
	}

	if insertedRoom.Price != room.Price {
		t.Errorf("expected Price %.2f, got %.2f", room.Price, insertedRoom.Price)
	}

	if insertedRoom.HotelID != insertedHotel.ID {
		t.Errorf("expected HotelID %v, got %v", insertedHotel.ID, insertedRoom.HotelID)
	}

	// Verify the hotel now contains the room ID
	updatedHotel, err := tdb.hotelStore.GetHotelByID(context.TODO(), insertedHotel.ID)
	if err != nil {
		t.Fatalf("error getting updated hotel: %v", err)
	}

	roomFound := false
	for _, id := range updatedHotel.Rooms {
		if id == insertedRoom.ID {
			roomFound = true
			break
		}
	}

	if !roomFound {
		t.Errorf("expected hotel to contain the room ID")
	}
}

// TestMongoRoomStore_GetRooms tests fetching rooms with optional filters
func TestMongoRoomStore_GetRooms(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping MongoDB integration test in short mode")
	}

	tdb := setupRoomTest(t)
	defer tdb.teardown(t)

	// First create a hotel to associate the rooms with
	hotel := &types.Hotel{
		Name:     "Test Hotel",
		Location: "Test Location",
		Rooms:    []primitive.ObjectID{},
		Rating:   4,
	}

	insertedHotel, err := tdb.hotelStore.Insert(context.TODO(), hotel)
	if err != nil {
		t.Fatalf("error inserting hotel: %v", err)
	}

	// Insert multiple test rooms with different attributes
	rooms := []*types.Room{
		{
			Size:    "small",
			Seaside: false,
			Price:   89.99,
			HotelID: insertedHotel.ID,
		},
		{
			Size:    "large",
			Seaside: true,
			Price:   149.99,
			HotelID: insertedHotel.ID,
		},
	}

	for _, room := range rooms {
		_, err := tdb.roomStore.InsertRoom(context.TODO(), room)
		if err != nil {
			t.Fatalf("error inserting room: %v", err)
		}
	}

	// Test 1: Get all rooms for the hotel
	hotelFilter := bson.M{"hotelID": insertedHotel.ID}
	fetchedRooms, err := tdb.roomStore.GetRooms(context.TODO(), hotelFilter)
	if err != nil {
		t.Fatalf("error getting rooms: %v", err)
	}

	if len(fetchedRooms) != len(rooms) {
		t.Fatalf("expected %d rooms, got %d", len(rooms), len(fetchedRooms))
	}

	// Test 2: Filter by seaside and price
	seasideFilter := bson.M{"hotelID": insertedHotel.ID, "seaside": true, "price": bson.M{"$gt": 100}}
	filteredRooms, err := tdb.roomStore.GetRooms(context.TODO(), seasideFilter)
	if err != nil {
		t.Fatalf("error getting rooms with filter: %v", err)
	}

	if len(filteredRooms) != 1 {
		t.Fatalf("expected 1 seaside room with price > 100, got %d", len(filteredRooms))
	}

	if !filteredRooms[0].Seaside {
		t.Errorf("expected seaside room, got non-seaside")
	}

	if filteredRooms[0].Price <= 100 {
		t.Errorf("expected room price > 100, got %.2f", filteredRooms[0].Price)
	}
} 