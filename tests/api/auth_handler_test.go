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

// Test login request parameters
type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Test login response structure
type loginResp struct {
	Token string `json:"token"`
	User  types.User `json:"user"`
}

// setupAuth creates a test server with auth routes configured
func setupAuth(t *testing.T) (*fiber.App, *mongo.Client, func()) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Initialize user store and handlers
	userStore := db.NewMongoUserStore(client)
	authHandler := api.NewAuthHandler(userStore)

	// Clean up the user collection before test
	_, err = client.Database(db.DBNAME).Collection("users").DeleteMany(context.TODO(), struct{}{})
	if err != nil {
		t.Fatalf("Error cleaning users collection: %v", err)
	}

	// Setup Fiber app
	app := fiber.New()
	
	// Setup auth routes manually
	app.Post("/api/auth/login", authHandler.HandleAuthentication)

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

// createTestUser creates a user for testing authentication
func createTestUser(t *testing.T, client *mongo.Client) *types.User {
	userStore := db.NewMongoUserStore(client)
	
	// Use a fixed ID for consistent test results
	userID := primitive.NewObjectID()
	
	// Create user params with plaintext password
	params := types.CreateUserParams{
		FirstName: "Test",
		LastName:  "User",
		Email:     "testauth@example.com",
		Password:  "password123",
	}
	
	// Create the user
	user, err := types.NewUserFromParams(params)
	if err != nil {
		t.Fatalf("Error creating user from params: %v", err)
	}
	
	// Set the ID explicitly
	user.ID = userID
	
	// Insert into database
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatalf("Error inserting test user: %v", err)
	}
	
	return user
}

// TestUserLogin tests the user login authentication
func TestUserLogin(t *testing.T) {
	app, client, cleanup := setupAuth(t)
	defer cleanup()
	
	// Create a test user first
	user := createTestUser(t, client)
	
	// Prepare login request
	login := loginReq{
		Email:    user.Email,
		Password: "password123", // Same as in createTestUser
	}
	
	body, err := json.Marshal(login)
	if err != nil {
		t.Fatalf("Error marshaling login request: %v", err)
	}
	
	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
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
	var loginResponse loginResp
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	
	// Verify token is present
	if loginResponse.Token == "" {
		t.Errorf("Expected token to be present in response")
	}
	
	// Verify user info
	if loginResponse.User.ID != user.ID {
		t.Errorf("Expected user ID %v, got %v", user.ID, loginResponse.User.ID)
	}
	
	if loginResponse.User.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, loginResponse.User.Email)
	}
}

// TestFailedLogin tests authentication failure due to wrong password
func TestFailedLogin(t *testing.T) {
	app, client, cleanup := setupAuth(t)
	defer cleanup()
	
	// Create a test user first
	user := createTestUser(t, client)
	
	// Prepare login request with wrong password
	login := loginReq{
		Email:    user.Email,
		Password: "wrongpassword",
	}
	
	body, err := json.Marshal(login)
	if err != nil {
		t.Fatalf("Error marshaling login request: %v", err)
	}
	
	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	
	// Execute request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()
	
	// Check status code - should be unauthorized or similar error
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Expected non-OK status code for wrong password, got %v", resp.StatusCode)
	}
	
	// Don't try to parse the error response as JSON, just check that it's not empty
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	errorBody := buf.String()
	
	if errorBody == "" {
		t.Errorf("Expected error message in response body")
	}
} 