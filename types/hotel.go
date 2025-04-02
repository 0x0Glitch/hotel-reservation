package types

import "go.mongodb.org/mongo-driver/bson/primitive"

// RoomType defines the available types of rooms in the hotel
// Using iota for auto-incrementing integer constants
type RoomType int

const(
	_ RoomType = iota      // Skip the first value (0)
	SingleRoomType         // Value: 1 - A room for one person
	DoubleRoomType         // Value: 2 - A room for two people
	SeaSideRoomType        // Value: 3 - A room with a sea view
	DeluxRoomType          // Value: 4 - A luxury room with premium amenities
)

// Hotel represents a hotel in the reservation system
// Contains basic information about the hotel and references to its rooms
type Hotel struct{
	ID 		 primitive.ObjectID     `bson:"_id,omitempty" json:"id"`        // Unique identifier for the hotel
	Name 	 string 		        `bson:"name" json:"name"`               // Name of the hotel
	Location string 	            `bson:"location" json:"location"`       // Physical location/address of the hotel
	Rooms 	 []primitive.ObjectID	`bson:"rooms" json:"rooms"`             // List of room IDs belonging to this hotel
	Rating 	 int					`bson:"rating" json:"rating"`           // Hotel rating (e.g., 1-5 stars)
}

// Room represents an individual room in a hotel
// Contains details about the room's features and pricing
type Room struct{
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"` // Unique identifier for the room
	Seaside	  bool 				     `bson:"seaside" json:"seaside"`           // Whether the room has a sea view
	Size 	  string			     `bson: "size" json:"size"`               // Size of the room (e.g., "large", "small")
	Price 	  float64			     `bson:"price",json:"price"`              // Cost per night in the room
	HotelID   primitive.ObjectID     `bson:"hotelID" json:"hotelID"`           // ID of the hotel this room belongs to
}
