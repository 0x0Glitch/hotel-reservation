package db

import (
	"context"

	"github.com/0x0Glitch/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
const usesrColl = "users"
type UserStore interface{
	GetUserById(context.Context,string) (*types.User,error)
}

type MongoUserStore struct{
	client *mongo.Client
	dbname string
	coll *mongo.Collection
}
func NewMongoUserStore(client *mongo.Client) *MongoUserStore{
	
	
	return &MongoUserStore{
		client: client,
		coll :client.Database(DBNAME).Collection(usesrColl),
	}

}
func (s *MongoUserStore) GetUserById(ctx context.Context,id string) (*types.User,error){
	var user types.User
	oid ,err := primitive.ObjectIDFromHex(id)
	if err != nil{
		return nil,err
	}

	if err := s.coll.FindOne(ctx,bson.M{"_id": oid}).Decode(&user); err != nil{
		return nil,err
	}
	return &user,nil
}
