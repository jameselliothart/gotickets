package tickets

import (
	"fmt"
	"sync"

	"github.com/jameselliothart/gotickets/cqrs"
)

type TicketCommandHandler struct {
	EventHandlers []cqrs.EventHandler
	CommandHandler cqrs.CommandHandler
}

func (h *TicketCommandHandler) Dispatch(cmd cqrs.Command) {
	ch := make(chan error)
	wgHandlers := &sync.WaitGroup{}
	wg := &sync.WaitGroup{}
	events, err := h.CommandHandler.HandleCommand(cmd)
	if err != nil {
		cqrs.LogWithCorrelation(cmd, fmt.Sprintf("%T: %v", h, err))
		return
	}
	for _, event := range events {
		for _, handler := range h.EventHandlers {
			wgHandlers.Add(1)
			go func(event cqrs.Event, handler cqrs.EventHandler, wg *sync.WaitGroup, ch chan<- error) {
				ch <- handler.HandleEvent(event)
				wg.Done()
			}(event, handler, wgHandlers, ch)
		}
		wg.Add(1)
		go func(event cqrs.Event, wg *sync.WaitGroup, ch <-chan error) {
			for err := range ch {
				if err != nil {
					cqrs.LogWithCorrelation(event, fmt.Sprintf("Error handling event %#v: %v", event, err))
				}
			}
			wg.Done()
		}(event, wg, ch)
	}
	wgHandlers.Wait()
	close(ch)
	wg.Wait()
	return
}
