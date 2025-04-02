package main

import (
	"context"
	"flag"
	"log"

	"github.com/0x0Glitch/hotel-reservation/api"
	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/middleware"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Fiber configuration for custom error handling
// This ensures all errors are returned in a consistent JSON format
var config = fiber.Config{
    // Override default error handler to return JSON instead of plain text
    ErrorHandler: func(c *fiber.Ctx, err error) error {
        return c.JSON(map[string]string{"error":err.Error()})
    },
}

// main is the entry point of the application
// It sets up the database connection, handlers, and starts the HTTP server
func main(){
	// Parse command line flags
	// You can specify a different port using: go run main.go -listenAddr=:8080
	listenAddr := flag.String("listenAddr",":5001","The listen address of the API server")
	flag.Parse()

	// Connect to MongoDB
	// The URI is defined in the db package
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil{
		log.Fatal(err)
	}
	
	// Initialize database stores
	// These provide access to different collections in MongoDB
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client,hotelStore)
	userStore := db.NewMongoUserStore(client)
	
	// Create a central store with all sub-stores
	store := &db.Store{
		Hotel: hotelStore,
		Room: roomStore,
		User: userStore,
	}
	
	// Initialize API handlers
	// These handle HTTP requests and use the stores to interact with the database
	userHandler := api.NewUserHandler(userStore)
	hotelHandler := api.NewHotelHandler(store)
	authHandler := api.NewAuthHandler(userStore)
	
	// Create a new Fiber app with our custom config
	app := fiber.New(config)
	
	// Create API routes
	// auth group is for non-authenticated endpoints
	auth := app.Group("/api")
	
	// apiv1 group requires JWT authentication for all routes
	apiv1 := app.Group("/api/v1",middleware.JWTAuthentication)	

	// Authentication routes
	// These don't require authentication to access
	auth.Post("/auth",authHandler.HandleAuthentication)

	// User routes
	// All of these require authentication
	apiv1.Post("/user",userHandler.HandlePostUser)         // Create a new user
	apiv1.Delete("/user/:id",userHandler.HandleDeleteUser) // Delete a user
	apiv1.Get("/user", userHandler.HandleGetUsers)         // Get all users
	apiv1.Get("/user/:id",userHandler.HandleGetUser)       // Get a specific user
	apiv1.Put("user/:id",userHandler.HandlePutUser)        // Update a user
	
	// Hotel routes
	// All of these require authentication
	apiv1.Get("/hotel",hotelHandler.HandleGetHotels)         // Get all hotels
	apiv1.Get("/hotel/:id",hotelHandler.HandleGetHotel)      // Get a specific hotel
	apiv1.Get("/hotel/:id/rooms",hotelHandler.HandleGetRooms) // Get rooms for a hotel

	// Start the server
	app.Listen(*listenAddr)
}


