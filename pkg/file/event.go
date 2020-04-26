package file

import (
	"time"
)

type EventTypeMask int

const (
	EventTypeCreate = EventTypeMask(1 << iota)
	EventTypeUpdate
	EventTypeDelete
)

type Event interface {
	Path() Path
	Timestamp() time.Time
	TypeMask() EventTypeMask
}
