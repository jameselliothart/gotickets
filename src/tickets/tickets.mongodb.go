package tickets

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type TicketsCollection struct {
	*mongo.Collection
}

func NewTicketsCollection(db *mongo.Database) *TicketsCollection {
	ticketsCollection := db.Collection("tickets")
	return &TicketsCollection{ticketsCollection}
}

func (mc *TicketsCollection) GetTickets() []Ticket {
	return []Ticket{
		{
			ID:      uuid.New().String(),
			Summary: "ticket 1",
		},
		{
			ID:      uuid.New().String(),
			Summary: "ticket 2",
		},
	}
}
