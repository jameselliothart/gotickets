package cqrs

import (
	"log"
	"time"

	"github.com/google/uuid"
)

type Correlatable interface {
	CorrelationID() uuid.UUID
}

type Command interface {
	Correlatable
}

type CommandDispatcher interface {
	Dispatch(Command)
}

type CommandHandler interface {
	HandleCommand(Command) ([]Event, error)
}

type Event interface {
	Correlatable
	CreatedAt() time.Time
}

type EventHandler interface {
	HandleEvent(Event) error
}

func LogWithCorrelation(c Correlatable, v ...interface{}) {
	log.Println(c.CorrelationID(), v)
}

type EventLogger struct{}

func (l EventLogger) HandleEvent(event Event) error {
	LogWithCorrelation(event, "Event created:", event)
	return nil
}
