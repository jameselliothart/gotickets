package tickets

import "github.com/lithammer/shortuuid/v3"

type TicketID struct {
	ID string `json:"id"`
}

type Ticket struct {
	TicketID
	Summary  string `json:"summary"`
}

type CreateTicketDto struct {
	Summary string
}

func NewTicketID() TicketID {
	return TicketID{ID: shortuuid.New()}
}

func NewTicket(summary string) Ticket {
	return Ticket{
		TicketID: NewTicketID(),
		Summary: summary,
	}
}

type DataHandler interface {
	GetTickets() []Ticket
	CreateTicket(CreateTicketDto) (TicketID, error)
}
