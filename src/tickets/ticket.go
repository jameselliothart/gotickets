package tickets

type Ticket struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
}

type DataHandler interface {
	GetTickets() []Ticket
}
