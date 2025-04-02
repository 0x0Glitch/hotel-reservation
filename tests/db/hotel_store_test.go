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

// hotelTestDB wraps the hotel store and client for testing
type hotelTestDB struct {
	client     *mongo.Client
	hotelStore db.HotelStore
}

// setupHotelTest initializes the test environment with a MongoDB connection
func setupHotelTest(t *testing.T) *hotelTestDB {
	// Connect to test database
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDBURI))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	
	// Clean up any previous test data to ensure a fresh start
	_, err = client.Database(db.DBNAME).Collection("hotels").DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		t.Fatalf("Error cleaning up test collection: %v", err)
	}

	return &hotelTestDB{
		client:     client,
		hotelStore: hotelStore,
	}
}

// teardown cleans up after tests
func (tdb *hotelTestDB) teardown(t *testing.T) {
	// Delete all test data
	_, err := tdb.client.Database(db.DBNAME).Collection("hotels").DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		t.Fatalf("Error cleaning up test collection: %v", err)
	}

	// Disconnect from MongoDB
	if err := tdb.client.Disconnect(context.TODO()); err != nil {
		t.Fatalf("Error disconnecting from MongoDB: %v", err)
	}
}

// TestMongoHotelStore_Insert tests inserting a new hotel
func TestMongoHotelStore_Insert(t *testing.T) {
	// Skip integration tests when running in short mode
	if testing.Short() {
		t.Skip("Skipping MongoDB integration test in short mode")
	}

	tdb := setupHotelTest(t)
	defer tdb.teardown(t)

	// Create a test hotel
	hotel := &types.Hotel{
		Name:     "Test Hotel",
		Location: "Test Location",
		Rooms:    []primitive.ObjectID{},
		Rating:   4,
	}

	// Insert the hotel
	insertedHotel, err := tdb.hotelStore.Insert(context.TODO(), hotel)
	if err != nil {
		t.Fatalf("error inserting hotel: %v", err)
	}

	// Verify hotel ID is set
	if insertedHotel.ID.IsZero() {
		t.Errorf("expected hotel ID to be set")
	}

	// Verify hotel fields match
	if insertedHotel.Name != hotel.Name {
		t.Errorf("expected Name %s, got %s", hotel.Name, insertedHotel.Name)
	}

	if insertedHotel.Location != hotel.Location {
		t.Errorf("expected Location %s, got %s", hotel.Location, insertedHotel.Location)
	}

	if insertedHotel.Rating != hotel.Rating {
		t.Errorf("expected Rating %d, got %d", hotel.Rating, insertedHotel.Rating)
	}
}

// TestMongoHotelStore_Update tests updating hotel fields
func TestMongoHotelStore_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping MongoDB integration test in short mode")
	}

	tdb := setupHotelTest(t)
	defer tdb.teardown(t)

	// Create and insert a test hotel
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

	// Update hotel's name - Note: MongoDB requires $set for updates
	filter := bson.M{"_id": insertedHotel.ID}
	update := bson.M{"$set": bson.M{"name": "Updated Hotel Name"}}

	err = tdb.hotelStore.Update(context.TODO(), filter, update)
	if err != nil {
		t.Fatalf("error updating hotel: %v", err)
	}

	// Verify the update worked
	updatedHotel, err := tdb.hotelStore.GetHotelByID(context.TODO(), insertedHotel.ID)
	if err != nil {
		t.Fatalf("error getting updated hotel: %v", err)
	}

	if updatedHotel.Name != "Updated Hotel Name" {
		t.Errorf("expected updated Name 'Updated Hotel Name', got '%s'", updatedHotel.Name)
	}
}

// TestMongoHotelStore_GetHotels tests fetching hotels with optional filters
func TestMongoHotelStore_GetHotels(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping MongoDB integration test in short mode")
	}

	tdb := setupHotelTest(t)
	defer tdb.teardown(t)

	// Insert multiple test hotels with different ratings
	hotels := []*types.Hotel{
		{
			Name:     "Hotel A",
			Location: "Location A",
			Rooms:    []primitive.ObjectID{},
			Rating:   3,
		},
		{
			Name:     "Hotel B",
			Location: "Location B",
			Rooms:    []primitive.ObjectID{},
			Rating:   5,
		},
	}

	for _, hotel := range hotels {
		_, err := tdb.hotelStore.Insert(context.TODO(), hotel)
		if err != nil {
			t.Fatalf("error inserting hotel: %v", err)
		}
	}

	// Test 1: Get all hotels
	fetchedHotels, err := tdb.hotelStore.GetHotels(context.TODO(), bson.M{})
	if err != nil {
		t.Fatalf("error getting hotels: %v", err)
	}

	if len(fetchedHotels) != len(hotels) {
		t.Fatalf("expected %d hotels, got %d", len(hotels), len(fetchedHotels))
	}

	// Test 2: Filter by rating
	ratingFilter := bson.M{"rating": 5}
	filteredHotels, err := tdb.hotelStore.GetHotels(context.TODO(), ratingFilter)
	if err != nil {
		t.Fatalf("error getting hotels with filter: %v", err)
	}

	if len(filteredHotels) != 1 {
		t.Fatalf("expected 1 hotel with rating 5, got %d", len(filteredHotels))
	}

	if filteredHotels[0].Rating != 5 {
		t.Errorf("expected hotel with rating 5, got %d", filteredHotels[0].Rating)
	}
}

// TestMongoHotelStore_GetHotelByID tests fetching a specific hotel by ID
func TestMongoHotelStore_GetHotelByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping MongoDB integration test in short mode")
	}

	tdb := setupHotelTest(t)
	defer tdb.teardown(t)

	// Create and insert a test hotel
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

	// Fetch the hotel by ID
	fetchedHotel, err := tdb.hotelStore.GetHotelByID(context.TODO(), insertedHotel.ID)
	if err != nil {
		t.Fatalf("error getting hotel by ID: %v", err)
	}

	// Verify hotel fields match
	if fetchedHotel.ID != insertedHotel.ID {
		t.Errorf("expected ID %v, got %v", insertedHotel.ID, fetchedHotel.ID)
	}

	if fetchedHotel.Name != insertedHotel.Name {
		t.Errorf("expected Name %s, got %s", insertedHotel.Name, fetchedHotel.Name)
	}

	if fetchedHotel.Location != insertedHotel.Location {
		t.Errorf("expected Location %s, got %s", insertedHotel.Location, fetchedHotel.Location)
	}

	if fetchedHotel.Rating != insertedHotel.Rating {
		t.Errorf("expected Rating %d, got %d", insertedHotel.Rating, fetchedHotel.Rating)
	}

	// Test with non-existent ID
	nonExistentID := primitive.NewObjectID()
	_, err = tdb.hotelStore.GetHotelByID(context.TODO(), nonExistentID)
	if err == nil {
		t.Errorf("expected error when getting non-existent hotel")
	}
} 