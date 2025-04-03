package db

import (
	"context"
	"testing"

	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Constants for the test database
const testDBName = "hotel-reservation-test"

// testDB wraps the user store for testing
type testDB struct {
	client    *mongo.Client
	userStore db.UserStore
}

// setup initializes the test database with a fresh MongoDB connection
func setup(t *testing.T) *testDB {
	// Skip tests if MongoDB is not available
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDBURI))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}

	return &testDB{
		client:    client,
		userStore: db.NewMongoUserStore(client),
	}
}

// teardown cleans up after tests by dropping the test collection
func (tdb *testDB) teardown(t *testing.T) {
	// Drop the collection to clean up test data
	if err := tdb.userStore.Drop(context.TODO()); err != nil {
		t.Fatalf("Error dropping test collection: %v", err)
	}

	// Disconnect from MongoDB
	if err := tdb.client.Disconnect(context.TODO()); err != nil {
		t.Fatalf("Error disconnecting from MongoDB: %v", err)
	}
}

// TestMongoUserStore_InsertUser tests inserting a user into MongoDB
func TestMongoUserStore_InsertUser(t *testing.T) {
	// Skip test if no MongoDB connection is available
	if testing.Short() {
		t.Skip("Skipping MongoDB integration test in short mode")
	}

	tdb := setup(t)
	defer tdb.teardown(t)

	// Create test user with explicit ID to avoid duplicate key issues
	user := &types.User{
		ID:                primitive.NewObjectID(),
		FirstName:         "John",
		LastName:          "Doe",
		Email:             "john@example.com",
		EncryptedPassword: "encrypted_password",
	}

	// Insert the user
	insertedUser, err := tdb.userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatalf("error inserting user: %v", err)
	}

	// Verify user ID is set
	if insertedUser.ID.IsZero() {
		t.Errorf("expected user ID to be set")
	}

	// Verify user fields match
	if insertedUser.FirstName != user.FirstName {
		t.Errorf("expected firstName %s, got %s", user.FirstName, insertedUser.FirstName)
	}

	if insertedUser.LastName != user.LastName {
		t.Errorf("expected lastName %s, got %s", user.LastName, insertedUser.LastName)
	}

	if insertedUser.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, insertedUser.Email)
	}
}

// TestMongoUserStore_GetUserById tests fetching a user by ID
func TestMongoUserStore_GetUserById(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping MongoDB integration test in short mode")
	}

	tdb := setup(t)
	defer tdb.teardown(t)

	// Create and insert a test user
	user := &types.User{
		ID:                primitive.NewObjectID(),
		FirstName:         "John",
		LastName:          "Doe",
		Email:             "john@example.com",
		EncryptedPassword: "encrypted_password",
	}

	insertedUser, err := tdb.userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatalf("error inserting user: %v", err)
	}

	// Fetch the user by ID
	fetchedUser, err := tdb.userStore.GetUserById(context.TODO(), insertedUser.ID.Hex())
	if err != nil {
		t.Fatalf("error getting user by ID: %v", err)
	}

	// Verify user fields match
	if fetchedUser.ID != insertedUser.ID {
		t.Errorf("expected ID %v, got %v", insertedUser.ID, fetchedUser.ID)
	}

	if fetchedUser.FirstName != insertedUser.FirstName {
		t.Errorf("expected firstName %s, got %s", insertedUser.FirstName, fetchedUser.FirstName)
	}

	if fetchedUser.LastName != insertedUser.LastName {
		t.Errorf("expected lastName %s, got %s", insertedUser.LastName, fetchedUser.LastName)
	}

	if fetchedUser.Email != insertedUser.Email {
		t.Errorf("expected email %s, got %s", insertedUser.Email, fetchedUser.Email)
	}

	// Test with non-existent ID
	nonExistentID := primitive.NewObjectID().Hex()
	_, err = tdb.userStore.GetUserById(context.TODO(), nonExistentID)
	if err == nil {
		t.Errorf("expected error when getting non-existent user")
	}
}

// TestMongoUserStore_GetUsers tests fetching all users
func TestMongoUserStore_GetUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping MongoDB integration test in short mode")
	}

	tdb := setup(t)
	defer tdb.teardown(t)

	// Insert multiple test users
	users := []*types.User{
		{
			ID:                primitive.NewObjectID(), // Explicit ID to avoid duplicates
			FirstName:         "John",
			LastName:          "Doe",
			Email:             "john@example.com",
			EncryptedPassword: "encrypted_password1",
		},
		{
			ID:                primitive.NewObjectID(), // Explicit ID to avoid duplicates
			FirstName:         "Jane",
			LastName:          "Smith",
			Email:             "jane@example.com",
			EncryptedPassword: "encrypted_password2",
		},
	}

	for _, user := range users {
		_, err := tdb.userStore.InsertUser(context.TODO(), user)
		if err != nil {
			t.Fatalf("error inserting user: %v", err)
		}
	}

	// Fetch all users
	fetchedUsers, err := tdb.userStore.GetUsers(context.TODO())
	if err != nil {
		t.Fatalf("error getting users: %v", err)
	}

	// Verify count matches
	if len(fetchedUsers) != len(users) {
		t.Fatalf("expected %d users, got %d", len(users), len(fetchedUsers))
	}

	// Check that all users are in the result set
	emails := make(map[string]bool)
	for _, user := range fetchedUsers {
		emails[user.Email] = true
	}

	for _, user := range users {
		if !emails[user.Email] {
			t.Errorf("expected to find user with email %s in retrieved users", user.Email)
		}
	}
}

// TestMongoUserStore_DeleteUser tests deleting a user
func TestMongoUserStore_DeleteUser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping MongoDB integration test in short mode")
	}

	tdb := setup(t)
	defer tdb.teardown(t)

	// Create and insert a test user
	user := &types.User{
		ID:                primitive.NewObjectID(),
		FirstName:         "John",
		LastName:          "Doe",
		Email:             "john@example.com",
		EncryptedPassword: "encrypted_password",
	}

	insertedUser, err := tdb.userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatalf("error inserting user: %v", err)
	}

	// Delete the user
	err = tdb.userStore.DeleteUser(context.TODO(), insertedUser.ID.Hex())
	if err != nil {
		t.Fatalf("error deleting user: %v", err)
	}

	// Verify the user is deleted by trying to fetch it
	_, err = tdb.userStore.GetUserById(context.TODO(), insertedUser.ID.Hex())
	if err == nil {
		t.Errorf("expected error after deleting user, got nil")
	}
}

// TestMongoUserStore_UpdateUser tests updating user fields
func TestMongoUserStore_UpdateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping MongoDB integration test in short mode")
	}

	tdb := setup(t)
	defer tdb.teardown(t)

	// Create and insert a test user
	user := &types.User{
		ID:                primitive.NewObjectID(),
		FirstName:         "John",
		LastName:          "Doe",
		Email:             "john@example.com",
		EncryptedPassword: "encrypted_password",
	}

	insertedUser, err := tdb.userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatalf("error inserting user: %v", err)
	}

	// Update user's first name - Note: $set is required for MongoDB updates
	filter := bson.M{"_id": insertedUser.ID}
	update := bson.M{"firstName": "JohnUpdated"}

	err = tdb.userStore.UpdateUser(context.TODO(), filter, update)
	if err != nil {
		t.Fatalf("error updating user: %v", err)
	}

	// Verify the update worked
	updatedUser, err := tdb.userStore.GetUserById(context.TODO(), insertedUser.ID.Hex())
	if err != nil {
		t.Fatalf("error getting updated user: %v", err)
	}

	if updatedUser.FirstName != "JohnUpdated" {
		t.Errorf("expected updated firstName 'JohnUpdated', got '%s'", updatedUser.FirstName)
	}
}

// TestMongoUserStore_GetUserByEmail tests fetching a user by email
func TestMongoUserStore_GetUserByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping MongoDB integration test in short mode")
	}

	tdb := setup(t)
	defer tdb.teardown(t)

	// Create and insert a test user
	email := "john@example.com"
	user := &types.User{
		ID:                primitive.NewObjectID(),
		FirstName:         "John",
		LastName:          "Doe",
		Email:             email,
		EncryptedPassword: "encrypted_password",
	}

	_, err := tdb.userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatalf("error inserting user: %v", err)
	}

	// Fetch the user by email
	fetchedUser, err := tdb.userStore.GetUserByEmail(context.TODO(), email)
	if err != nil {
		t.Fatalf("error getting user by email: %v", err)
	}

	// Verify email matches
	if fetchedUser.Email != email {
		t.Errorf("expected email %s, got %s", email, fetchedUser.Email)
	}

	// Test with non-existent email
	_, err = tdb.userStore.GetUserByEmail(context.TODO(), "nonexistent@example.com")
	if err == nil {
		t.Errorf("expected error when getting user with non-existent email")
	}
} 