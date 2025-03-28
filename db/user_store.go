package db

import (
	"context"
	"fmt"

	"github.com/0x0Glitch/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
const usesrColl = "users"
type UserStore interface{

	GetUserById(context.Context,string) (*types.User,error)
	GetUsers(context.Context) ([]*types.User,error)
	InsertUser(context.Context,*types.User) (*types.User,error)
	DeleteUser(context.Context, string) error
	UpdateUser(ctx context.Context,filter bson.M,update bson.M) error
	Drop(context.Context) error


}

type MongoUserStore struct{
	client *mongo.Client
	dbname string
	coll *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client,dbname string) *MongoUserStore{
	
	
	return &MongoUserStore{
		client: client,
		coll :client.Database(dbname).Collection(usesrColl),
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

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User,error){
	var users []*types.User
	cur,err := s.coll.Find(ctx,bson.M{})
	if err != nil{
		return nil,err
	}
	
	
	if err := cur.All(ctx,&users); err!=nil{
		return nil,err
	}

	return users,nil
}

func (s *MongoUserStore) InsertUser(ctx context.Context,user *types.User)(*types.User,error) {
	res, err := s.coll.InsertOne(ctx,user)
	if err != nil{
		return nil,err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user,nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context,id string) error {
	oid ,err := primitive.ObjectIDFromHex(id)
	if err != nil{
		return err
	}
	
	_, err = s.coll.DeleteOne(ctx,bson.M{"_id":oid})
	if err != nil{
		return err
	}
	return nil
	
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M,update bson.M)error{
	_,err := s.coll.UpdateOne(ctx, filter,bson.M{"$set": update})
	if err != nil{
		return err
	}
	return nil
}

func (s *MongoUserStore) Drop(ctx context.Context) error{
	fmt.Println("----dropping")
	return s.coll.Drop(ctx)
}






