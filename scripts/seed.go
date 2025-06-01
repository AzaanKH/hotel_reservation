package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/AzaanKH/hotel_reservation/api"
	"github.com/AzaanKH/hotel_reservation/db"
	"github.com/AzaanKH/hotel_reservation/db/fixtures"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found or cannot be loaded: %v", err)
	}
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(db.DBURI))

	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Booking: db.NewMongoBookingStore(client),
		Hotel:   hotelStore,
	}
	user := fixtures.AddUser(store, "Buddy", "Hield", false)
	fmt.Println("Buddy Hield ->", api.CreateTokenFromUser(user))
	adminUser := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin ->", api.CreateTokenFromUser(adminUser))
	hotel := fixtures.AddHotel(store, "Buddys Hotel", "Pakistan", 4, nil)
	room := fixtures.AddRoom(store, "small", true, 69.99, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("booking ->", booking.ID)

	for i := range 100 {
		name := fmt.Sprintf("random hotel name %d", i)
		location := fmt.Sprintf("location%d", i)

		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}
	// // Properly close the connection
	if err := client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
}
