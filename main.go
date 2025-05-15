package main

import (
	"context"
	"flag"
	"log"

	"github.com/AzaanKH/hotel_reservation/api"
	"github.com/AzaanKH/hotel_reservation/db"
	"github.com/AzaanKH/hotel_reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userColl = "users"

var config = fiber.Config{

	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	listtenAddr := flag.String("listenAddr", ":5001", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(db.DBURI))

	if err != nil {
		log.Fatal(err)
	}

	user := types.User{
		FirstName: "Buddy",
		LastName:  "Hield",
	}
	client.Database(db.DBURI).Collection(userColl).InsertOne(context.TODO(), user)

	// handlers init
	var (
		userStore   = db.NewMongoUserStore(client)
		userHandler = api.NewUserHandler(userStore)
		hotelStore  = db.NewMongoHotelStore(client)
		roomStore   = db.NewMongoRoomStore(client, hotelStore)
		store       = &db.Store{
			User:  userStore,
			Room:  roomStore,
			Hotel: hotelStore,
		}
		hotelHandler = api.NewHotelHandler(store)
		app          = fiber.New(config)
		apiv1        = app.Group("api/v1")
	)

	// user handlers
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("hotel/:id/rooms", hotelHandler.HandleGetRooms)
	app.Listen(*listtenAddr)

}
