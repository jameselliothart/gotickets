package tickets

import (
	"time"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v3"
)

type TicketID struct {
	ID string `json:"id"`
}

func NewTicketID() TicketID {
	return TicketID{ID: shortuuid.New()}
}

type TicketFields struct {
	Summary string `json:"summary"`
}

type Ticket struct {
	TicketID
	TicketFields
}

func NewTicket(summary string) Ticket {
	return Ticket{
		TicketID: NewTicketID(),
		TicketFields: TicketFields{
			Summary: summary,
		},
	}
}

type CreateTicketDto struct {
	Summary string
}

type CreateTicketCmd struct {
	corrID      uuid.UUID
	TicketFields
}

func NewCreateTicketCmd(summary string) CreateTicketCmd {
	return CreateTicketCmd{
		corrID: uuid.New(),
		TicketFields: TicketFields{
			Summary: summary,
		},
	}
}

func (c CreateTicketCmd) CorrelationID() uuid.UUID {
	return c.corrID
}

type TicketCreatedEvent struct {
	corrID uuid.UUID
	createdAt time.Time
	TicketFields
}

func NewTicketCreatedEvent(cmd CreateTicketCmd) TicketCreatedEvent {
	return TicketCreatedEvent{
		corrID: cmd.CorrelationID(),
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

type DataHandler interface {
	GetTickets() []Ticket
}
