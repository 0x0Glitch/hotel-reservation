package api

import (
	"errors"

	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles HTTP requests related to user operations
// It uses the UserStore to interact with the database
type UserHandler struct{
	userStore db.UserStore // Database interface for user operations
}

// NewUserHandler creates a new UserHandler with the provided UserStore
// Factory function to create handlers with dependency injection
func NewUserHandler(userStore db.UserStore) *UserHandler{
	return &UserHandler{
		userStore: userStore,
	}
}

// HandleGetUser processes requests to get a single user by ID
// GET /api/users/:id
func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error{
	// Extract user ID from URL parameters
	var id = c.Params("id")
	
	// Fetch the user from the database
	user, err:= h.userStore.GetUserById(c.Context(), id)
	if err != nil{
		// If user is not found, return a friendly error message
		if errors.Is(err, mongo.ErrNoDocuments){
			return c.JSON(map[string]string{"error":"not found"})
		}
		// For any other error, return it directly
		return err
	}
	// Return the user as JSON
	return c.JSON(user)
}

// HandleGetUsers processes requests to get all users
// GET /api/users
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error{
	// Fetch all users from the database
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil{
		return err
	}
	// Return users as JSON array
	return c.JSON(users)
}

// HandlePostUser processes requests to create a new user
// POST /api/users
func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	// Create variable to hold request data
	var params types.CreateUserParams
	
	// Parse request body into CreateUserParams struct
	if err := c.BodyParser(&params); err != nil{
		return err
	}
	
	// Validate the user input
	if errors := params.Validate(); len(errors) > 0{
		// If validation fails, return errors to the client
		return c.JSON(errors)
	}
	
	// Create new user from params (includes password hashing)
	user, err := types.NewUserFromParams(params)
	if err != nil{
		return err
	}
	
	// Save the user to the database
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil{
		return err
	}
	
	// Return the newly created user with ID
	return c.JSON(insertedUser)
}

// HandleDeleteUser processes requests to delete a user
// DELETE /api/users/:id
func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error{
	// Extract user ID from URL parameters
	userID := c.Params("id")
	
	// Delete the user from the database
	if err := h.userStore.DeleteUser(c.Context(), userID); err != nil{
		return err
	}
	
	// Return success message with deleted ID
	return c.JSON(map[string]string{"deleted": userID})
}

// HandlePutUser processes requests to update a user
// PUT /api/users/:id
func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	// Create a map to hold the update fields
	var update bson.M
	
	// Parse the JSON body into the update map
	if err := c.BodyParser(&update); err != nil {
		return err
	}

	// Extract user ID from URL parameters
	userID := c.Params("id")
	
	// Convert string ID to MongoDB ObjectID
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	// Create filter to find the user by ID
	filter := bson.M{"_id": oid}
	
	// Update the user in the database
	if err := h.userStore.UpdateUser(c.Context(), filter, update); err != nil {
		return err
	}

	// Return success message with updated ID
	return c.JSON(map[string]string{"updated": userID})
}

