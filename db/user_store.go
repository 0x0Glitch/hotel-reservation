package db

import (
	"context"
	"fmt"

	"github.com/0x0Glitch/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Name of the MongoDB collection for users
const usesrColl = "users"

// UserStore defines the interface for user data operations
// Any implementation of UserStore must provide these methods
type UserStore interface{
	GetUserByEmail(context.Context,string) (*types.User,error)   // Find a user by email address
	GetUserById(context.Context,string) (*types.User,error)      // Find a user by their ID
	GetUsers(context.Context) ([]*types.User,error)              // Get all users
	InsertUser(context.Context,*types.User) (*types.User,error)  // Add a new user
	DeleteUser(context.Context, string) error                    // Remove a user
	UpdateUser(ctx context.Context,filter bson.M,update bson.M) error // Update user information
	Drop(context.Context) error                                  // Drop the entire users collection (dangerous!)
}

// MongoUserStore implements the UserStore interface with MongoDB
// It handles all user-related database operations
type MongoUserStore struct{
	client *mongo.Client              // MongoDB client connection
	dbname string                     // Database name
	coll *mongo.Collection            // Reference to the users collection
}

// NewMongoUserStore creates a new MongoUserStore with the provided MongoDB client
// This is a factory function that sets up the connection to the users collection
func NewMongoUserStore(client *mongo.Client) *MongoUserStore{
	return &MongoUserStore{
		client: client,
		coll: client.Database(DBNAME).Collection(usesrColl),
	}
}

// GetUserById retrieves a user by their ID
// It takes a context and ID string, and returns the user or an error
func (s *MongoUserStore) GetUserById(ctx context.Context,id string) (*types.User,error){
	var user types.User
	
	// Convert string ID to MongoDB ObjectID
	oid ,err := primitive.ObjectIDFromHex(id)
	if err != nil{
		return nil,err
	}

	// Find the user document and decode it into the user variable
	if err := s.coll.FindOne(ctx,bson.M{"_id": oid}).Decode(&user); err != nil{
		return nil,err
	}
	return &user,nil
}

// GetUsers retrieves all users from the database
// Returns a slice of user pointers or an error
func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User,error){
	var users []*types.User
	
	// Find all user documents
	cur,err := s.coll.Find(ctx,bson.M{})
	if err != nil{
		return nil,err
	}
	
	// Decode all results into the users slice
	if err := cur.All(ctx,&users); err!=nil{
		return nil,err
	}

	return users,nil
}

// InsertUser adds a new user to the database
// Takes a user object and returns the inserted user with ID or an error
func (s *MongoUserStore) InsertUser(ctx context.Context,user *types.User)(*types.User,error) {
	// Insert the user document
	res, err := s.coll.InsertOne(ctx,user)
	if err != nil{
		return nil,err
	}
	
	// Update the user object with the generated ID
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user,nil
}

// DeleteUser removes a user from the database by ID
// Returns an error if the operation failed
func (s *MongoUserStore) DeleteUser(ctx context.Context,id string) error {
	// Convert string ID to MongoDB ObjectID
	oid ,err := primitive.ObjectIDFromHex(id)
	if err != nil{
		return err
	}
	
	// Delete the user document
	_, err = s.coll.DeleteOne(ctx,bson.M{"_id":oid})
	if err != nil{
		return err
	}
	return nil
}

// UpdateUser modifies user information
// Takes a filter to select the user and an update document
func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M,update bson.M)error{
	// Update the user document using $set to avoid replacing the entire document
	_,err := s.coll.UpdateOne(ctx, filter,bson.M{"$set": update})
	if err != nil{
		return err
	}
	return nil
}

// Drop deletes the entire users collection
// This is typically used for testing or resetting the database
// WARNING: This will delete ALL users and cannot be undone
func (s *MongoUserStore) Drop(ctx context.Context) error{
	fmt.Println("----dropping")
	return s.coll.Drop(ctx)
}

// GetUserByEmail finds a user by their email address
// Email addresses are unique in the system
func (s *MongoUserStore) GetUserByEmail(ctx context.Context,email string) (*types.User,error){
	var user types.User
	// Find and decode the user document
	if err := s.coll.FindOne(ctx,bson.M{"email": email}).Decode(&user); err != nil{
		return nil,err
	}
	return &user,nil
}





