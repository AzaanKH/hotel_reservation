package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AzaanKH/hotel_reservation/db"
	"github.com/AzaanKH/hotel_reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99.9,
		},
		{
			Size:  "normal",
			Price: 199.9,
		},
		{
			Size:  "king",
			Price: 345.9,
		},
	}
	insteredHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insteredHotel)
	for _, room := range rooms {
		room.HotelID = insteredHotel.ID
		_, err := roomStore.Insert(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func main() {
	seedHotel("Buddys Hotel", "Pakistan", 4)
	seedHotel("Buddys Hotel #2", "Lahore", 4)
	seedHotel("Buddys Hotel #3", "Bahamas", 5)
}

func init() {
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(db.DBURI))

	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Collection("hotels").Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
