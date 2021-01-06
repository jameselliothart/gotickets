package tickets

import (
	"time"

	"github.com/google/uuid"
)

type TicketCreatedEvent struct {
	corrID    uuid.UUID
	createdAt time.Time
	TicketFields
}

func NewTicketCreatedEvent(cmd CreateTicketCmd) TicketCreatedEvent {
	return TicketCreatedEvent{
		corrID:    cmd.CorrelationID(),
		createdAt: time.Now().UTC(),
		TicketFields: TicketFields{
			Summary: cmd.Summary,
		},
	}
}

func (e TicketCreatedEvent) CorrelationID() uuid.UUID {
	return e.corrID
}

func (e TicketCreatedEvent) CreatedAt() time.Time {
	return e.createdAt
}

type TicketCompletedEvent struct {
	corrID    uuid.UUID
	createdAt time.Time
	Ticket    *Ticket
}

func NewTicketCompletedEvent(cmd CompleteTicketCmd, ticket *Ticket) TicketCompletedEvent {
	return TicketCompletedEvent{
		corrID:    cmd.CorrelationID(),
		createdAt: time.Now().UTC(),
		Ticket:    ticket,
	}
}

func (e TicketCompletedEvent) CorrelationID() uuid.UUID {
	return e.corrID
}

func (e TicketCompletedEvent) CreatedAt() time.Time {
	return e.createdAt
}
