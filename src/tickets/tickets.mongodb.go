package tickets

import (
	"context"
	"fmt"
	"log"
	"sync"

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

func (c *TicketsCollection) getTickets() []Ticket {
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
	return TicketID{ID: ticket.ID}, nil
}

func (c *TicketsCollection) HandleEvent(event cqrs.Event, wg *sync.WaitGroup, ch chan<- error) {
	var err error
	switch e := event.(type) {
	case TicketCreatedEvent:
		ticket := NewTicket(e.Summary)
		id, err := c.create(ticket)
		if err == nil {
			cqrs.LogWithCorrelation(event, fmt.Sprintf("Created ticket id: '%v'", id))
		}
	default:
		err = errors.Errorf("%T does not recognize event: %#v", c, e)
	}
	ch <- err
	wg.Done()
}

func (c *TicketsCollection) HandleQuery(query interface{}) (tickets []Ticket) {
	switch query.(type) {
	case ActiveTicketsQuery:
		tickets = c.getTickets()
	default:
		log.Printf("Query not recognized: %#v", c)
	}
	return
}
