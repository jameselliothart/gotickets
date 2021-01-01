package tickets

import (
	"log"

	"github.com/jameselliothart/gotickets/cqrs"
)

type TicketCommandHandler struct {
	Handlers []cqrs.EventHandler
}

func (h *TicketCommandHandler) HandleCommand(cmd cqrs.Command) {
	switch c := cmd.(type) {
	case CreateTicketCmd:
		ticketCreated := NewTicketCreatedEvent(c)
		for _, handler := range h.Handlers {
			if err := handler.HandleEvent(ticketCreated); err != nil {
				log.Printf("Error handling event %#v: %v", ticketCreated, err)
			}
		}
	default:
		log.Printf("Command not recognized: %#v", c)
	}
	return
}
