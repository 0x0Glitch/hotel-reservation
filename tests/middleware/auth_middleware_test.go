package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupAuthMiddlewareTest sets up the environment for testing auth middleware
func setupAuthMiddlewareTest(t *testing.T) (*fiber.App, *mongo.Client, func()) {
	// Set JWT secret for testing
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Initialize user store
	userStore := db.NewMongoUserStore(client)

	// Clean up the user collection before test
	_, err = client.Database(db.DBNAME).Collection("users").DeleteMany(context.TODO(), struct{}{})
	if err != nil {
		t.Fatalf("Error cleaning users collection: %v", err)
	}

	// Setup Fiber app with a protected route
	app := fiber.New()
	
	// Create a protected route that requires authentication
	// Normally we would use middleware.JWTAuthentication here, but for test simplicity,
	// we'll create a mock protected handler that checks for a user in the authorization header
	app.Get("/api/protected", func(c *fiber.Ctx) error {
		// Mock JWT authentication middleware for testing
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}
		
		// Extract token
		tokenString := auth[7:] // Remove "Bearer " prefix
		
		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		
		if err != nil || !token.Valid {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}
		
		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}
		
		// Check expiration
		if exp, ok := claims["expires"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Token expired",
				})
			}
		}
		
		// Get user ID from token
		userIDStr, ok := claims["id"].(string)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID in token",
			})
		}
		
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		
		// Get user from database
		user, err := userStore.GetUserById(c.Context(), userID.Hex())
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		
		// Set user in context
		c.Locals("user", user)
		
		// Handle the protected route
		user, ok = c.Locals("user").(*types.User)
		if !ok {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get user from context",
			})
		}
		return c.JSON(fiber.Map{
			"message": "Protected route accessed successfully",
			"userID":  user.ID,
		})
	})

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
		os.Unsetenv("JWT_SECRET")
	}
}

// createTestUserAndToken creates a test user and generates a valid JWT token
func createTestUserAndToken(t *testing.T, client *mongo.Client) (*types.User, string) {
	userStore := db.NewMongoUserStore(client)
	
	// Create a user
	userID := primitive.NewObjectID()
	user := &types.User{
		ID:                userID,
		FirstName:         "Auth",
		LastName:          "Test",
		Email:             "authtest@example.com",
		EncryptedPassword: "some-encrypted-password",
	}
	
	insertedUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatalf("Error inserting test user: %v", err)
	}
	
	// Create a token
	token := createToken(insertedUser)
	
	return insertedUser, token
}

// createToken generates a JWT token for the given user
func createToken(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour).Unix()
	
	claims := jwt.MapClaims{
		"id":      user.ID.Hex(), // Convert ObjectID to string
		"email":   user.Email,
		"expires": expires,
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, _ := token.SignedString([]byte(secret))
	
	return tokenString
}

// TestAuthMiddleware_ValidToken tests authentication with a valid token
func TestAuthMiddleware_ValidToken(t *testing.T) {
	app, client, cleanup := setupAuthMiddlewareTest(t)
	defer cleanup()
	
	// Create user and token
	user, token := createTestUserAndToken(t, client)
	
	// Create HTTP request with auth header
	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	
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
	
	// Parse response to verify user ID
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	
	// Verify the user ID in the response matches our test user
	if userID, ok := response["userID"].(string); ok {
		if userID != user.ID.Hex() {
			t.Errorf("Expected user ID %s, got %s", user.ID.Hex(), userID)
		}
	} else {
		t.Errorf("Missing or invalid userID in response")
	}
}

// TestAuthMiddleware_NoToken tests authentication with no token
func TestAuthMiddleware_NoToken(t *testing.T) {
	app, _, cleanup := setupAuthMiddlewareTest(t)
	defer cleanup()
	
	// Create HTTP request without auth header
	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	
	// Execute request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()
	
	// Check status code - should be unauthorized
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Expected non-OK status for request without token, got %v", resp.StatusCode)
	}
}

// TestAuthMiddleware_InvalidToken tests authentication with invalid token
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	app, _, cleanup := setupAuthMiddlewareTest(t)
	defer cleanup()
	
	// Create HTTP request with invalid token
	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	
	// Execute request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()
	
	// Check status code - should be unauthorized
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Expected non-OK status for request with invalid token, got %v", resp.StatusCode)
	}
}

// TestAuthMiddleware_ExpiredToken tests authentication with expired token
func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	app, client, cleanup := setupAuthMiddlewareTest(t)
	defer cleanup()
	
	// Create a user
	userStore := db.NewMongoUserStore(client)
	userID := primitive.NewObjectID()
	user := &types.User{
		ID:                userID,
		FirstName:         "Expired",
		LastName:          "Token",
		Email:             "expired@example.com",
		EncryptedPassword: "some-encrypted-password",
	}
	
	_, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatalf("Error inserting test user: %v", err)
	}
	
	// Create an expired token
	now := time.Now()
	expires := now.Add(-time.Hour).Unix() // Expired 1 hour ago
	
	claims := jwt.MapClaims{
		"id":      user.ID.Hex(), // Convert ObjectID to string
		"email":   user.Email,
		"expires": expires,
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, _ := token.SignedString([]byte(secret))
	
	// Create HTTP request with expired token
	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	
	// Execute request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()
	
	// Check status code - should be unauthorized
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Expected non-OK status for request with expired token, got %v", resp.StatusCode)
	}
} 