package main

import (
	"context"
	"flag"
	"log"

	"github.com/AzaanKH/hotel_reservation/api"
	"github.com/AzaanKH/hotel_reservation/api/middleware"
	"github.com/AzaanKH/hotel_reservation/db"
	"github.com/AzaanKH/hotel_reservation/types"
	"github.com/gofiber/fiber/v2"

	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found or cannot be loaded: %v", err)
	}
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
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			User:    userStore,
			Room:    roomStore,
			Hotel:   hotelStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("/api")
		apiv1          = app.Group("api/v1", middleware.JWTAuthentication(userStore))
		admin          = apiv1.Group("/admin", middleware.AdminAuth)
	)
	// auth handler
	auth.Post("/auth", authHandler.HandleAuthenticate)

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
	// room handlers
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	// cancel booking

	// booking handler
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)

	// admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)
	app.Listen(*listtenAddr)

}
