package types

import (
	"fmt"

	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Constants for user-related validations
const (
	bcryptCost = 12              // Security level for password hashing (higher = more secure but slower)
	miniFirstNameLen = 2         // Minimum allowed length for first name
	miniLastNameLen= 2           // Minimum allowed length for last name
	miniPasswordLen = 7          // Minimum allowed length for password
)

// UpdateUserParams defines the data needed to update a user
// This is used when updating an existing user's information
type UpdateUserParams struct{
	FirstName   string `json:"firstName"` // User's first name
	LastName 	string `json:"lastName"`  // User's last name
}

// CreateUserParams defines the data needed to create a new user
// This is used during user registration
type CreateUserParams struct{
	FirstName   string `json:"firstName"` // User's first name
	LastName 	string `json:"lastName"`  // User's last name
	Email 		string `json:"email"`     // User's email address (must be valid format)
	Password 	string `json:"password"`  // User's password (plain text during creation only)
}

// User represents a user in the hotel reservation system
// This is the main user model stored in the database
type User struct {
    ID                primitive.ObjectID `bson:"_id" json:"id"`                      // Unique identifier (MongoDB ObjectID)
    FirstName         string             `bson:"firstName" json:"firstName"`         // User's first name
    LastName          string             `bson:"lastName"  json:"lastName"`          // User's last name
    Email             string             `bson:"email"     json:"email"`             // User's email address
    EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`         // Password hash (not sent in JSON responses)
}

// NewUserFromParams creates a new User object from the provided parameters
// It handles password encryption and generates a new unique ID
func NewUserFromParams(params CreateUserParams) (*User,error){
	// Generate a secure hash of the password
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil{
		return nil,err
	}
	// Create and return a new user with a unique ID
	return &User{FirstName: params.FirstName,
		LastName: params.LastName,
		Email: params.Email,
		EncryptedPassword: string(encpw),
		ID:  primitive.NewObjectID(),
	},nil
}

// Validate checks if the CreateUserParams contains valid data
// Returns a map of field names to error messages for any invalid fields
func (params CreateUserParams) Validate() map[string]string{
	errors := map[string]string{}
	
	// Check first name length
	if len(params.FirstName)<miniFirstNameLen{
		errors["firstName"] = fmt.Sprintf("firstName length should be at least %d characters",miniFirstNameLen)
	}
	
	// Check last name length
	if len(params.LastName)<miniFirstNameLen{
		 errors["lastName"] = fmt.Sprintf("LastNamelength should be at least %d characters",miniLastNameLen)
	}

	// Check password length
	if len(params.Password)<miniPasswordLen{
		errors["password"]=fmt.Sprintf("minimum password length should be at least %d characters",miniPasswordLen)
	}
	
	// Check email format
	if !IsEmailValid(params.Email){
		errors["email"] = fmt.Sprintf("Email is invalid")
	}
	return errors
}

// IsEmailValid checks if the provided email string has a valid format
// Uses regex pattern matching to validate email format
func IsEmailValid(e string) bool{
		var emailRegex = regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
		return emailRegex.MatchString(e)
	}

// IsValidPassword checks if the provided plain text password matches the encrypted password
// Used during login to verify user credentials
func IsValidPassword(encpw,pw string)bool{
		return bcrypt.CompareHashAndPassword([]byte(encpw),[]byte(pw))==nil
}