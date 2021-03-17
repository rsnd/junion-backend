package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/rsnd/junion-backend/config"
)

// Ctx is a wait timer
var Ctx context.Context

// Client is a mongo db cluster
var Client *mongo.Client

// UserCollection is a mongo db collection for the users model
var UserCollection *mongo.Collection

// EmailverificationsCollection is a mongo db collection for the emailverifications model
var EmailverificationsCollection *mongo.Collection

// EventsCollection is a mongo db collection for the events model
var EventsCollection *mongo.Collection

// ConversationsCollection is a mongo db collection for the conversations model
var ConversationsCollection *mongo.Collection

// PollsCollection is a mongo db collection for the polls model
var PollsCollection *mongo.Collection

// ConnectDB connects to a specified database provided by
// the environment database URL
func ConnectDB() {
	currentConfig := config.GetConfig()

	Client, err := mongo.NewClient(options.Client().ApplyURI(currentConfig["DATABASE_URL"]))
	if err != nil {
		log.Fatal(err)
	}
	Ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = Client.Connect(Ctx)
	if err != nil {
		log.Fatal(err)
	}
	// defer Client.Disconnect(Ctx)
	err = Client.Ping(Ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	db := Client.Database(currentConfig["DATABASE_NAME"])
	if err != nil {
		log.Fatal(err)
	}
	UserCollection = db.Collection("users")
	EmailverificationsCollection = db.Collection("emailverifications")
	EventsCollection = db.Collection("events")
	ConversationsCollection = db.Collection("conversations")
	PollsCollection = db.Collection("polls")

	fmt.Println("Database successfully connected and pinged.")
}
