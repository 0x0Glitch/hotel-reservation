package db

import (
	"context"

	"github.com/0x0Glitch/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RoomStore defines the interface for room data operations
// Any implementation of RoomStore must provide these methods
type RoomStore interface{
	InsertRoom(context.Context,*types.Room) (*types.Room, error)  // Add a new room
	GetRooms(context.Context,bson.M)([]*types.Room,error)         // Get rooms with optional filters
}

// MongoRoomStore implements the RoomStore interface with MongoDB
// It handles all room-related database operations
type MongoRoomStore struct{
	client *mongo.Client           // MongoDB client connection
	coll *mongo.Collection         // Reference to the rooms collection
	HotelStore                     // Embedded HotelStore for hotel operations
}

// NewMongoRoomStore creates a new MongoRoomStore with the provided MongoDB client
// This is a factory function that sets up the connection to the rooms collection
// It requires a HotelStore because rooms belong to hotels and need to update them
func NewMongoRoomStore(client *mongo.Client,hotelStore HotelStore) *MongoRoomStore{
	return &MongoRoomStore{
		client: client,
		coll: client.Database(DBNAME).Collection("rooms"),
		HotelStore: hotelStore,
	}
}

// InsertRoom adds a new room to the database
// It also updates the associated hotel to include the new room ID
func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room)(*types.Room,error){
	// Insert the room document
	resp,err :=s.coll.InsertOne(ctx,room)
	if err != nil{
		return nil,err
	}
	
	// Update the room object with the generated ID
	room.ID = resp.InsertedID.(primitive.ObjectID)
	
	// Create filter and update to add the room ID to the hotel's rooms array
	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}

	// Update the hotel document to include this room
	if err := s.HotelStore.Update(ctx,filter,update);err != nil{
		return nil,err
	}
	
	return room,nil
}

// GetRooms retrieves rooms from the database
// The filter parameter allows for querying specific rooms (e.g., by hotel ID)
func (s *MongoRoomStore) GetRooms(ctx context.Context,filter bson.M) ([]*types.Room,error){
	// Find rooms matching the filter
	resp,err := s.coll.Find(ctx,filter)
	if err != nil{
		return nil,err
	}
	
	// Decode all results into the rooms slice
	var rooms []*types.Room
	if err := resp.All(ctx,&rooms);err != nil{
		return nil,err
	}
	return rooms,nil
}