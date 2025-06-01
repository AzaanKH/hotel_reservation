package api

import (
	"context"
	"log"
	"testing"

	"github.com/AzaanKH/hotel_reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

const (
	testdburi = "mongodb://localhost:27017"
	dbname    = "hotel-reservation-test"
)

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(testdburi))

	if err != nil {
		log.Fatal(err)
	}
	hotel := db.NewMongoHotelStore(client)
	return &testdb{
		Store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Hotel:   hotel,
			Room:    db.NewMongoRoomStore(client, hotel),
			Booking: db.NewMongoBookingStore(client),
		},
		client: client,
	}
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(dbname).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
