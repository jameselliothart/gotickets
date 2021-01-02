package tickets

import (
	"fmt"
	"sync"

	"github.com/jameselliothart/gotickets/cqrs"
)

type TicketCommandHandler struct {
	Handlers []cqrs.EventHandler
}

func (h *TicketCommandHandler) HandleCommand(cmd cqrs.Command) {
	ch := make(chan error)
	wgHandlers := &sync.WaitGroup{}
	wg := &sync.WaitGroup{}
	switch c := cmd.(type) {
	case CreateTicketCmd:
		ticketCreated := NewTicketCreatedEvent(c)
		for _, handler := range h.Handlers {
			wgHandlers.Add(1)
			go handler.HandleEvent(ticketCreated, wgHandlers, ch)
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, ch <-chan error) {
			for err := range ch {
				if err != nil {
					cqrs.LogWithCorrelation(ticketCreated, fmt.Sprintf("Error handling event %#v: %v", ticketCreated, err))
				}
			}
			wg.Done()
		}(wg, ch)
	default:
		cqrs.LogWithCorrelation(cmd, fmt.Sprintf("Command not recognized: %#v", c))
	}
	wgHandlers.Wait()
	close(ch)
	wg.Wait()
	return
}
