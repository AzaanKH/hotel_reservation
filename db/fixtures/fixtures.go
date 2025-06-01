package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AzaanKH/hotel_reservation/db"
	"github.com/AzaanKH/hotel_reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fname, lname string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fname, lname),
		FirstName: fname,
		LastName:  lname,
		Password:  fmt.Sprintf("%s_%s", fname, lname),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.Admin = admin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddHotel(store *db.Store, name, loc string, rating int, room []primitive.ObjectID) *types.Hotel {
	var roomIDS = room
	if room == nil {
		roomIDS = []primitive.ObjectID{}
	}

	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomIDS,
		Rating:   rating,
	}

	insteredHotel, err := store.Hotel.Insert(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insteredHotel
}

func AddRoom(store *db.Store, size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := store.Room.Insert(context.TODO(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBooking(store *db.Store, userID, roomID primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}
	booking, err := store.Booking.InsertBooking(context.TODO(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return booking
}
