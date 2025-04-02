package main

import (
	"context"
	"log"

	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var hotelStore db.HotelStore
var roomStore db.RoomStore
var ctx = context.Background()
var userStore db.UserStore

func seedHotel(rating int,name,location string) {

	hotel := types.Hotel{
		Name: name,
		Location: location,
		Rooms: []primitive.ObjectID{},
		Rating: rating,
	}
	rooms := []types.Room{
		{
		Size: "small",
		Price: 99,
	},{
		Size: "normal",
		Price: 899,
	},

	}
	insertedhotel,err := hotelStore.Insert(ctx,&hotel)
	if err != nil{
		log.Fatal(err)
	}
	

	for _,room := range rooms{
		room.HotelID = insertedhotel.ID
		_,err := roomStore.InsertRoom(ctx,&room)
		if err != nil{
		log.Fatal(err)
	}
	}
}
func main(){
	
	seedHotel(3,"Bellucia","France")
	seedHotel(4,"Sandrosso","Roorkee")
	seedUser("anshuman","yadav","anshumaniitre9@gmail.com")
}

func init(){
	var err error
	client, err := mongo.Connect(context.TODO(),options.Client().ApplyURI(db.DBURI))
	if err != nil{
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err!= nil{
		log.Fatal(err)
	}
	
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client,hotelStore)
	userStore = db.NewMongoUserStore(client)
	_ = roomStore
	
}
func seedUser(fname,lname,email string){
	user,err := types.NewUserFromParams(types.CreateUserParams{
		Email: email,
		FirstName: fname,
		LastName: lname,
		Password: "supersecurepassword",
		
	})
	if err != nil{
		log.Fatal(err)
	}
	_, err = userStore.InsertUser(ctx,user)
	if err != nil{
		log.Fatal(err)
	}

}



