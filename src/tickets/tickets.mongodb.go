package tickets

import (
	"context"
	"fmt"
	"log"

	"github.com/jameselliothart/gotickets/cqrs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TicketsCollection struct {
	*mongo.Collection
}

func NewTicketsCollection(db *mongo.Database) *TicketsCollection {
	ticketsCollection := db.Collection("tickets")
	return &TicketsCollection{ticketsCollection}
}

func (c *TicketsCollection) getTickets(filter interface{}) []Ticket {
	tickets := []Ticket{}
	cursor, err := c.Find(context.TODO(), filter)
	if err == nil {
		cursor.All(context.TODO(), &tickets)
		log.Printf("Retrieved %v tickets", len(tickets))
	} else {
		log.Printf("Error getting all tickets: %v\n", err)
	}
	return tickets
}

func (tc *TicketsCollection) getTicket(ticketID string) (ticket Ticket, err error) {
	filter := bson.D{primitive.E{Key: "ticketid.id", Value: ticketID}}
	err = tc.FindOne(context.TODO(), filter).Decode(&ticket)
	return
}

func (tc *TicketsCollection) create(ticket Ticket) (TicketID, error) {
	_, err := tc.InsertOne(context.TODO(), ticket)
	if err != nil {
		return TicketID{}, err
	}
	return TicketID{ID: ticket.ID}, nil
}

func (tc *TicketsCollection) updateTicket(ticket Ticket) error {
	filter := bson.D{primitive.E{Key: "ticketid.id", Value: ticket.ID}}
	result := tc.FindOneAndReplace(context.TODO(), filter, ticket)
	return result.Err()
}

func (tc *TicketsCollection) HandleEvent(event cqrs.Event) error {
	var err error
	switch e := event.(type) {
	case TicketCreatedEvent:
		ticket := NewTicket(e.Summary)
		id, err := tc.create(ticket)
		if err == nil {
			cqrs.LogWithCorrelation(e, fmt.Sprintf("Created ticket id: '%v'", id))
		}
	case TicketCompletedEvent:
		err := tc.updateTicket(*e.Ticket)
		if err == nil {
			cqrs.LogWithCorrelation(e, fmt.Sprintf("Completed ticket id: '%s'", e.Ticket.ID))
		}
	default:
		err = fmt.Errorf("%T does not recognize event: %#v", tc, e)
	}
	return err
}

func (tc *TicketsCollection) HandleQuery(query interface{}) (tickets []Ticket) {
	switch query.(type) {
	case OpenTicketsQuery:
		closed := bson.D{primitive.E{Key: "$in", Value: bson.A{ClosedStatuses()}}}
		notClosed := bson.D{primitive.E{Key: "$not", Value: closed}}
		statusOpen := bson.D{primitive.E{Key: "status", Value: notClosed}}
		tickets = tc.getTickets(statusOpen)
	default:
		log.Printf("Query not recognized: %#v", tc)
	}
	return
}

func (tc *TicketsCollection) HandleCommand(cmd cqrs.Command) ([]cqrs.Event, error) {
	switch c := cmd.(type) {
	case CreateTicketCmd:
		ticketCreated := NewTicketCreatedEvent(c)
		return []cqrs.Event{ticketCreated}, nil
	case CompleteTicketCmd:
		ticket, err := tc.getTicket(c.ID)
		if err != nil {
			return []cqrs.Event{}, err
		}
		if err = ticket.Complete(); err != nil {
			return []cqrs.Event{}, err
		}
		return []cqrs.Event{NewTicketCompletedEvent(c, &ticket)}, nil
	default:
		return []cqrs.Event{}, fmt.Errorf("%T does not recognize command: %#v", tc, c)
	}
}
