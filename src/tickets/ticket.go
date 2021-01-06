package tickets

import (
	"fmt"

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

type Status string

const (
	Active    Status = "active"
	Completed        = "completed"
	Cancelled        = "cancelled"
)

type Ticket struct {
	TicketID
	TicketFields
	Status Status
}

func NewTicket(summary string) Ticket {
	return Ticket{
		TicketID: NewTicketID(),
		TicketFields: TicketFields{
			Summary: summary,
		},
		Status: Active,
	}
}

func ClosedStatuses() []string {
	return []string{Completed, Cancelled}
}

func (t *Ticket) Complete() error {
	closedTicketError := fmt.Errorf("Ticket '%s' is already closed: status '%s'", t.ID, t.Status)
	switch t.Status {
	case "":
		t.Status = Completed
		return nil
	case Active:
		t.Status = Completed
		return nil
	case Completed:
		return closedTicketError
	case Cancelled:
		return closedTicketError
	default:
		return fmt.Errorf("Ticket '%#v' has unrecognized status: '%s'", t, t.Status)
	}
}

type CreateTicketDto struct {
	Summary string
}
