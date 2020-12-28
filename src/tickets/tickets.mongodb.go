package tickets

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TicketsCollection struct {
	*mongo.Collection
}

func NewTicketsCollection(db *mongo.Database) *TicketsCollection {
	ticketsCollection := db.Collection("tickets")
	return &TicketsCollection{ticketsCollection}
}

func (c *TicketsCollection) GetTickets() []Ticket {
	tickets := []Ticket{}
	cursor, err := c.Find(context.TODO(), bson.M{})
	if err == nil {
		cursor.All(context.TODO(), &tickets)
		log.Printf("Retrieved %v documents", len(tickets))
	} else {
		log.Printf("Error getting all tickets: %v\n", err)
	}
	return tickets
}

func (c *TicketsCollection) CreateTicket(dto CreateTicketDto) (TicketID, error) {
	ticket := NewTicket(dto.Summary)
	result, err := c.InsertOne(context.TODO(), ticket)
	if err != nil {
		return TicketID{}, err
	}
	log.Printf("Created ticket with mongo id: '%v', ticket id: '%v'", result.InsertedID, ticket.TicketID)
	return ticket.TicketID, err
}