package infrastructure

import (
	"context"
	"html/template"
	"log"
	"time"

	"github.com/jameselliothart/gotickets/cqrs"
	"github.com/jameselliothart/gotickets/tickets"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ticketsController tickets.TicketsController

func Startup() {
	layout := baseLayout()
	ticketsController.RegisterTemplates(layout)
	ticketsController.RegisterRoutes()

	ticketsCollection := tickets.NewTicketsCollection(connectToMongo("mongodb://localhost:27017"))
	ticketsController.QueryHandler = ticketsCollection
	ticketsController.CommandDispatcher = &tickets.TicketCommandHandler{
		CommandHandler: ticketsCollection,
		EventHandlers: []cqrs.EventHandler{
			ticketsCollection,
			new(cqrs.EventLogger),
		},
	}
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

func baseLayout() *template.Template {
	const basePath = "templates"
	layout := template.Must(template.ParseFiles(basePath + "/_layout.html"))
	template.Must(layout.ParseFiles(basePath+"/_header.html", basePath+"/_footer.html"))
	return layout
}
