package main

import (
	"context"
	"flag"
	"log"

	"github.com/0x0Glitch/hotel-reservation/api"
	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
const dburi = "mongodb://localhost:27017/"
const dbname = "hotel-reservation"
const userColl = "users"
// Create a new fiber instance with custom config
var config = fiber.Config{
    // Override default error handler
    ErrorHandler: func(c *fiber.Ctx, err error) error {
        return c.JSON(map[string]string{"error":err.Error()})
    },
}

// ...
func main(){
	listenAddr := flag.String("listenAddr",":5001","The listen address of the API server")
	flag.Parse()

	
	
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil{
		log.Fatal(err)
	}
	
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	
	
	app := fiber.New(config)

	apiv1 := app.Group("/api/v1")	
	apiv1.Get("/users", userHandler.HandleGetUsers)

	apiv1.Get("/user/:id",userHandler.HandleGetUser)

	app.Listen(*listenAddr)


}


