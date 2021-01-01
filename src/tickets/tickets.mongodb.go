package tickets

import (
	"context"
	"fmt"
	"log"

	"github.com/jameselliothart/gotickets/cqrs"
	"github.com/pkg/errors"
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

func (c *TicketsCollection) create(ticket Ticket) (TicketID, error) {
	_, err := c.InsertOne(context.TODO(), ticket)
	if err != nil {
		return TicketID{}, err
	}
	// log.Printf("Created ticket with: MongoID '%v' | TicketID '%v'", result.InsertedID, ticket.TicketID)
	return TicketID{ID: ticket.ID}, nil
}

func (c *TicketsCollection) HandleEvent(event cqrs.Event) error {
	switch e := event.(type) {
	case TicketCreatedEvent:
		ticket := NewTicket(e.Summary)
		id, err := c.create(ticket)
		if err == nil {
			cqrs.LogWithCorrelation(event, fmt.Sprintf("Created ticket id: '%v'", id))
		}
		return err
	default:
		return errors.Errorf("%T does not recognize event: %#v", c, e)
	}
}