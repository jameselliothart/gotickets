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
	// Database Config
	clientOptions := options.Client().ApplyURI(conn)
	client, err := mongo.NewClient(clientOptions)

	//Set up a context required by mongo.Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	//To close the connection at the end
	defer cancel()

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	}
	log.Println("Connected to mongo")
	db := client.Database("gotickets")
	return db
}
