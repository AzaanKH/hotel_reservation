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
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(admin bool, fname, lname, email, password string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.Admin = admin
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
}

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
	seedUser(false, "Buddy", "Hield", "buddybuckets24@mail.com", "buddybuckets2424")
	seedUser(true, "admin", "admin", "admin@mail.com", "buddybuckets24245728582")

	// users, err := userStore.GetUsers(context.Background())
	// if err != nil {
	// 	log.Fatal("Error retrieving users:", err)
	// }
	// fmt.Printf("Found %d users in database\n", len(users))
	// for _, u := range users {
	// 	fmt.Printf("User: %s %s (%s)\n", u.FirstName, u.LastName, u.Email)
	// }

	// time.Sleep(2 * time.Second)

	// // Properly close the connection
	if err := client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().
		ApplyURI(db.DBURI))

	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}
