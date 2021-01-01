package cqrs

import (
	"time"

	"github.com/google/uuid"
)

type Command interface {
	CorrelationID() uuid.UUID
}

type CommandHandler interface {
	HandleCmd(Command)
}

type Event interface {
	CorrelationID() uuid.UUID
	CreatedAt() time.Time
}

type EventHandler interface {
	HandleEvent(Event) error
}