package tickets

type TicketQueryHandler interface {
	HandleQuery(interface{}) []Ticket
}

type OpenTicketsQuery struct{}
