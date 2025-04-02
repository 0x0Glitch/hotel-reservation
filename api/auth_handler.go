package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthHandler handles HTTP requests related to authentication
// It processes login requests and generates JWT tokens
type AuthHandler struct{
	userStore db.UserStore // Database interface for user operations
}

// NewAuthHandler creates a new AuthHandler with the provided UserStore
// Factory function to create handlers with dependency injection
func NewAuthHandler(userStore db.UserStore) *AuthHandler{
	return &AuthHandler{
		userStore: userStore,
	}
}

// AuthParams defines the data needed for authentication
// Used for parsing login request bodies
type AuthParams struct{
	Email    string    `json:"email"`    // User's email address
	Password string    `json:"password"` // User's password (plain text for verification)
}

// AuthResponse defines the data returned after successful authentication
// Contains both the user data and the JWT token
type AuthResponse struct{
	User  *types.User `json:"User"`  // User information
	Token string      `json:"token"` // JWT token for authentication
}

// HandleAuthentication processes login requests
// POST /api/auth/login
func (h *AuthHandler) HandleAuthentication(c *fiber.Ctx) error{
	// Parse login parameters from request body
	var params AuthParams
	if err := c.BodyParser(&params); err != nil{
		return err
	}
	
	// Find user by email
	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil{
		// If user not found, return invalid credentials
		// This is a security best practice - don't reveal if the email exists
		if errors.Is(err, mongo.ErrNoDocuments){
			return fmt.Errorf("invalid credentials")
		}
		return err
	}
	
	// Verify password matches
	if !types.IsValidPassword(user.EncryptedPassword, params.Password){
		return fmt.Errorf("Invalid credentials")
	}
	
	// Generate JWT token for the user
	token := createTokenFromUser(user)
	
	// Create and return the response
	resp := AuthResponse{
		User:  user,
		Token: token,
	}
	return c.JSON(resp)
}

// createTokenFromUser generates a JWT token for the authenticated user
// The token contains user ID, email, and expiration time
func createTokenFromUser(user *types.User) string{
	now := time.Now()
	// Token expires in 4 hours
	expires := now.Add(time.Hour*4).Unix()
	
	// Create JWT claims (payload data)
	claims := jwt.MapClaims{
		"id":      user.ID,       // User ID for identification
		"email":   user.Email,    // Email for reference
		"expires": expires,       // Expiration timestamp
	}
	
	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Get the signing key from environment variables
	secret := os.Getenv("JWT_SECRET")
	
	// Sign the token with the secret key
	tokenStr, err := token.SignedString([]byte(secret))
	fmt.Println(secret)
	
	// Handle signing errors
	if err != nil{
		fmt.Println("Failed to sign token with secret", err)
	}
	
	return tokenStr
}