package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0x0Glitch/hotel-reservation/api"
	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Setup creates a test server with the specified MongoDB client and test database
func setup(t *testing.T) (*fiber.App, *mongo.Client, func()) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Initialize user store and handler
	userStore := db.NewMongoUserStore(client)
	userHandler := api.NewUserHandler(userStore)

	// Clean up the user collection before test
	_, err = client.Database(db.DBNAME).Collection("users").DeleteMany(context.TODO(), struct{}{})
	if err != nil {
		t.Fatalf("Error cleaning users collection: %v", err)
	}

	// Setup Fiber app
	app := fiber.New()
	
	// Setup user routes manually
	app.Post("/api/users", userHandler.HandlePostUser)
	app.Get("/api/users", userHandler.HandleGetUsers)
	app.Get("/api/users/:id", userHandler.HandleGetUser)
	app.Delete("/api/users/:id", userHandler.HandleDeleteUser)
	app.Put("/api/users/:id", userHandler.HandlePutUser)

	// Return test server and cleanup function
	return app, client, func() {
		// Clean up after test
		_, err := client.Database(db.DBNAME).Collection("users").DeleteMany(context.TODO(), struct{}{})
		if err != nil {
			t.Fatalf("Error cleaning users collection: %v", err)
		}
		if err := client.Disconnect(context.TODO()); err != nil {
			t.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}
}

// TestCreateUser tests the creation of a new user through the API
func TestCreateUser(t *testing.T) {
	app, _, cleanup := setup(t)
	defer cleanup()

	// Create test user
	userParams := types.CreateUserParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}

	// Convert to JSON payload
	body, err := json.Marshal(userParams)
	if err != nil {
		t.Fatalf("Error marshaling user params: %v", err)
	}

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

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
	var user types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify user fields
	if user.FirstName != userParams.FirstName {
		t.Errorf("Expected first name %s, got %s", userParams.FirstName, user.FirstName)
	}
	if user.LastName != userParams.LastName {
		t.Errorf("Expected last name %s, got %s", userParams.LastName, user.LastName)
	}
	if user.Email != userParams.Email {
		t.Errorf("Expected email %s, got %s", userParams.Email, user.Email)
	}
	if user.ID.IsZero() {
		t.Errorf("Expected user ID to be set")
	}
}

// TestGetUser tests retrieving a user by ID through the API
func TestGetUser(t *testing.T) {
	app, client, cleanup := setup(t)
	defer cleanup()

	// Create a user directly in the database first
	userStore := db.NewMongoUserStore(client)
	user := types.User{
		ID:                primitive.NewObjectID(),
		FirstName:         "Jane",
		LastName:          "Smith",
		Email:             "jane@example.com",
		EncryptedPassword: "somehashedpassword",
	}

	_, err := userStore.InsertUser(context.TODO(), &user)
	if err != nil {
		t.Fatalf("Error inserting test user: %v", err)
	}

	// Create HTTP request to get the user
	req := httptest.NewRequest(http.MethodGet, "/api/users/"+user.ID.Hex(), nil)

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
	var fetchedUser types.User
	if err := json.NewDecoder(resp.Body).Decode(&fetchedUser); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify user fields
	if fetchedUser.ID != user.ID {
		t.Errorf("Expected user ID %v, got %v", user.ID, fetchedUser.ID)
	}
	if fetchedUser.FirstName != user.FirstName {
		t.Errorf("Expected first name %s, got %s", user.FirstName, fetchedUser.FirstName)
	}
	if fetchedUser.LastName != user.LastName {
		t.Errorf("Expected last name %s, got %s", user.LastName, fetchedUser.LastName)
	}
	if fetchedUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, fetchedUser.Email)
	}
}

// TestGetUsers tests retrieving all users through the API
func TestGetUsers(t *testing.T) {
	app, client, cleanup := setup(t)
	defer cleanup()

	// Create multiple users directly in the database
	userStore := db.NewMongoUserStore(client)
	users := []types.User{
		{
			ID:                primitive.NewObjectID(),
			FirstName:         "User",
			LastName:          "One",
			Email:             "user1@example.com",
			EncryptedPassword: "hashedpw1",
		},
		{
			ID:                primitive.NewObjectID(),
			FirstName:         "User",
			LastName:          "Two",
			Email:             "user2@example.com",
			EncryptedPassword: "hashedpw2",
		},
	}

	for _, u := range users {
		user := u // Create a copy to avoid issues with loop variable capture
		_, err := userStore.InsertUser(context.TODO(), &user)
		if err != nil {
			t.Fatalf("Error inserting test user: %v", err)
		}
	}

	// Create HTTP request to get all users
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)

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
	var fetchedUsers []types.User
	if err := json.NewDecoder(resp.Body).Decode(&fetchedUsers); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify we got the expected number of users
	if len(fetchedUsers) != len(users) {
		t.Errorf("Expected %d users, got %d", len(users), len(fetchedUsers))
	}
} 