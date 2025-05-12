package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/AzaanKH/hotel_reservation/api"
	"github.com/AzaanKH/hotel_reservation/db"
	"github.com/AzaanKH/hotel_reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const userColl = "users"

var config = fiber.Config{
    
    ErrorHandler: func(c *fiber.Ctx, err error) error {
       return c.JSON(map[string]string{"error" : err.Error()})
    },
}

func main() {

	listtenAddr := flag.String("listenAddr", ":5001", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(dburi))

	if err != nil {
		log.Fatal(err)
	}

	user := types.User{
		FirstName: "Buddy",
		LastName:  "Hield",
	}
	client.Database(dburi).Collection(userColl).InsertOne(context.TODO(), user)

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	apiv1 := app.Group("api/v1")
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	app.Listen(*listtenAddr)
	fmt.Println()
}
