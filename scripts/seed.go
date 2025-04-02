package main

import (
	"context"
	"log"

	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global variables used throughout the seed script
var client *mongo.Client        // MongoDB client
var hotelStore db.HotelStore    // Interface to work with hotels collection
var roomStore db.RoomStore      // Interface to work with rooms collection
var ctx = context.Background()  // Context for database operations
var userStore db.UserStore      // Interface to work with users collection

// seedHotel creates a new hotel with the given parameters and adds two rooms to it
// This is a helper function to populate the database with sample hotel data
func seedHotel(rating int, name, location string) {
	// Create a new hotel object
	hotel := types.Hotel{
		Name: name,
		Location: location,
		Rooms: []primitive.ObjectID{},  // Empty list to be filled with room IDs
		Rating: rating,
	}
	
	// Define sample rooms to add to this hotel
	rooms := []types.Room{
		{
			Size: "small",
			Price: 99,
		}, {
			Size: "normal",
			Price: 899,
		},
	}
	
	// Insert the hotel into the database
	insertedhotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	
	// Insert each room and associate it with the hotel
	for _, room := range rooms {
		room.HotelID = insertedhotel.ID  // Set the hotel ID reference
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// main is the entry point for the seed script
// It calls the seeding functions to populate the database with initial data
func main() {
	// Seed sample hotels with rooms
	seedHotel(3, "Bellucia", "France")
	seedHotel(4, "Sandrosso", "Roorkee")
	
	// Seed a sample user
	seedUser("anshuman", "yadav", "anshumaniitre9@gmail.com")
}

// init is called before main() automatically by Go
// It sets up the database connection and initializes the stores
func init() {
	var err error
	
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	
	// Drop the existing database to start fresh
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	
	// Initialize stores for database operations
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}

// seedUser creates a new user with the given parameters
// This is a helper function to populate the database with sample user data
func seedUser(fname, lname, email string) {
	// Create a new user from parameters
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email: email,
		FirstName: fname,
		LastName: lname,
		Password: "supersecurepassword",  // Note: In a real app, use strong unique passwords
	})
	if err != nil {
		log.Fatal(err)
	}
	
	// Insert the user into the database
	_, err = userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}



