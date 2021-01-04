package tickets

import (
	"github.com/google/uuid"
	"github.com/jameselliothart/gotickets/cqrs"
	"github.com/pkg/errors"
)

type CreateTicketCmd struct {
	corrID uuid.UUID
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

type CompleteTicketCmd struct {
	corrID uuid.UUID
	TicketID
}

func NewCompleteTicketCmd(ticketID string) CompleteTicketCmd {
	return CompleteTicketCmd{
		corrID: uuid.New(),
		TicketID: TicketID{ID: ticketID},
	}
}

func NewStatusChangeCmd(status string, ticketID string) cqrs.Command {
	switch status {
	case "completed":
		return NewCompleteTicketCmd(ticketID)
	default:
		panic(errors.Errorf("Unrecognized status: '%s'", status))
	}
}

func (c CompleteTicketCmd) CorrelationID() uuid.UUID {
	return c.corrID
}