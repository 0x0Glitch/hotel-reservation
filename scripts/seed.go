package main

import (
	"context"
	"fmt"
	"log"

	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main(){
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(),options.Client().ApplyURI(db.DBURI))
	if err != nil{
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client,db.DBNAME)
	roomStore := db.NewRoomStore(client,db.DBNAME)
	hotel := types.Hotel{
		Name: "Sandrosso",
		Location: "Roorkee",
	}
	rooms := []types.Room{
		{
		Type: types.SingleRoomType,
		BasePrice: 999,
	},{
		Type: types.DeluxRoomType,
		BasePrice: 899,
	},

	}
	insertedhotel,err := hotelStore.InsertHotel(ctx,&hotel)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(insertedhotel)

	for _,room := range rooms{
		insertedRoom,err := roomStore.InsertRoom(ctx,&room)

		room.HotelID = insertedhotel.ID
		if err != nil{
		log.Fatal(err)
	}
	fmt.Println(insertedRoom)
	}
	

}



