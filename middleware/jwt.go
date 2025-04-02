package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthentication is a middleware function that checks if the request has a valid JWT token
// This ensures that only authenticated users can access protected routes
// It extracts the token from the X-Api-Token header, validates it, and checks if it's expired
func JWTAuthentication(c *fiber.Ctx) error{
	// Get the token from the request header
	token := c.Get("X-Api-Token")
	
	// Validate the token and get its claims (payload data)
	claims, err := validateToken(token)
	if err != nil {
		return err
	}
	
	// Check if the token has expired
	// First convert the expiration time from float64 to int64
	expiresFloat := claims["expires"].(float64)
	expires := int64(expiresFloat)
	
	// Compare current time with expiration time
	if time.Now().Unix() > expires {
		return fmt.Errorf("token expired")
	}
	
	// If token is valid and not expired, continue to the next middleware or handler
	c.Next()
	return nil
}

// validateToken checks if a JWT token is valid and returns its claims
// It verifies the token signature using the JWT_SECRET environment variable
// Returns the token claims if valid, or an error if not
func validateToken(tokenStr string) (jwt.MapClaims, error) {
	// Parse and validate the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token uses the correct signing method (HMAC in this case)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing methods", token.Header)
			return nil, fmt.Errorf("Unauthorized")
		}
		
		// Get the secret key from environment variables
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	
	// Handle parsing errors
	if err != nil {
		fmt.Println("failed to parse JWT token", err)
		return nil, fmt.Errorf("unauthorized")
	}
	
	// Check if the token is valid overall
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("unauthorized")
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	
	return claims, nil
}

