package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DBNAME     string
	DBURI      string
	TestDBNAME string
)

type Pagination struct {
	Limit int64
	Page  int64
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found or cannot be loaded: %v", err)
	}
	DBNAME = os.Getenv("MONGO_DB_NAME")
	DBURI = os.Getenv("MONGO_DB_URI")
	TestDBNAME = os.Getenv("MONGO_test_DBNAME")
}
