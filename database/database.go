package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Collection *mongo.Collection
var TeamCollection *mongo.Collection
var PendingCollection *mongo.Collection

func ConnectDB() {
	uri := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	db := client.Database("testRideDB")
	Collection = db.Collection("attendees")
	TeamCollection = db.Collection("team_members")
	PendingCollection = db.Collection("pending_attendees")

	log.Println("Connected to MongoDB ✅")
}
