package infrastructure

import (
	"context"
	"log"
	"time"

	"github.com/jameselliothart/gotickets/tickets"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ticketsController tickets.TicketsController

func Startup() {
	// ticketsController.TicketsTemplate = templates["tickets.html"]
	db := connectToMongo("mongodb://localhost:27017")
	ticketsController.DAL = tickets.NewTicketsCollection(db)
	ticketsController.RegisterRoutes()
}

func connectToMongo(conn string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(conn)
	client, err := mongo.NewClient(clientOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	defer cancel()

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	}
	log.Println("Connected to mongo")
	db := client.Database("gotickets")
	return db
}
