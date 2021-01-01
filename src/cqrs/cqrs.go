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

type CommandHandler interface {
	HandleCommand(Command)
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
