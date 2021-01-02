package cqrs

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Correlatable interface {
	CorrelationID() uuid.UUID
}

type Command interface {
	Correlatable
}

type CommandHandler interface {
	HandleCommand(Command)
}

type Event interface {
	Correlatable
	CreatedAt() time.Time
}

type EventHandler interface {
	HandleEvent(Event, *sync.WaitGroup, chan<- error)
}

func LogWithCorrelation(c Correlatable, v ...interface{}) {
	log.Println(c.CorrelationID(), v)
}

type EventLogger struct {}

func (l EventLogger) HandleEvent(event Event, wg *sync.WaitGroup, ch chan<- error) {
	LogWithCorrelation(event, "Event created:", event)
	ch <- nil
	wg.Done()
}
