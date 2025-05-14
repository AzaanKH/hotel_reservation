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

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(db.DBURI))

	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	hotel := types.Hotel{
		Name:     "Buddy",
		Location: "San Francisco",
		Rooms:    []primitive.ObjectID{},
	}

	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 99.9,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 199.9,
		},
		{
			Type:      types.SeaSideRoomType,
			BasePrice: 145.9,
		},
	}

	// room := types.Room{
	// 	Type:      types.SingleRoomType,
	// 	BasePrice: 99.99,
	// }
	insteredHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insteredHotel)
	for _, room := range rooms {
		room.HotelID = insteredHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}
}
