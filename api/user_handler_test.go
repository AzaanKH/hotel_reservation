package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AzaanKH/hotel_reservation/db"
	"github.com/AzaanKH/hotel_reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	db.UserStore
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
	return &testdb{
		UserStore: db.NewMongoUserStore(client, dbname),
	}
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)
	params := types.CreateUserParams{
		FirstName: "Buddy",
		LastName:  "Hied",
		Email:     "ab@mail.com",
		Password:  "avcw2341sa",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	fmt.Println(user)
	if user.FirstName != params.FirstName {
		t.Errorf("expected first name %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected last name %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
}

func TestDeleteUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)
	app.Delete("/:id", userHandler.HandleDeleteUser)
	app.Get("/:id", userHandler.HandleGetUser)
	params := types.CreateUserParams{
		FirstName: "Buddy",
		LastName:  "Hied",
		Email:     "ab@mail.com",
		Password:  "avcw2341sa",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	deleteURL := fmt.Sprintf("/%s", user.ID)
	deleteReq := httptest.NewRequest("DELETE", deleteURL, nil)
	deleteResp, err := app.Test(deleteReq)
	if err != nil {
		t.Error(err)
	}
	print(deleteResp.StatusCode)
	if deleteResp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK but got %d", deleteResp.StatusCode)
	}

	getURL := fmt.Sprintf("/%s", user.ID)
	getReq := httptest.NewRequest("GET", getURL, nil)
	getResp, err := app.Test(getReq)
	if err != nil {
		t.Error(err)
	}

	// Verify we get a not found status
	if getResp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status NotFound but got %d", getResp.StatusCode)
	}

	// 4. Verify error response
	var errorResp map[string]string
	json.NewDecoder(getResp.Body).Decode(&errorResp)

	if errorResp["error"] != "not found" {
		t.Errorf("expected error 'not found' but got %v", errorResp)
	}
}

// func TestUpdateUser(t *testing.T) {
// 	tdb := setup(t)
// 	defer tdb.teardown(t)
// 	app := fiber.New()
// 	userHandler := NewUserHandler(tdb.UserStore)
// 	app.Post("/", userHandler.HandlePostUser)
// 	params := types.CreateUserParams{
// 		FirstName: "Buddy",
// 		LastName:  "Hied",
// 		Email:     "ab@mail.com",
// 		Password:  "avcw2341sa",
// 	}
// 	b, _ := json.Marshal(params)

// 	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
// 	req.Header.Add("Content-Type", "application/json")
// 	resp, err := app.Test(req)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	var user types.User
// 	json.NewDecoder(resp.Body).Decode(&user)

// 	fmt.Println(user)
// 	if user.FirstName != params.FirstName {
// 		t.Errorf("expected first name %s but got %s", params.FirstName, user.FirstName)
// 	}
// 	if user.LastName != params.LastName {
// 		t.Errorf("expected last name %s but got %s", params.LastName, user.LastName)
// 	}
// 	if user.Email != params.Email {
// 		t.Errorf("expected email %s but got %s", params.Email, user.Email)
// 	}
// }
