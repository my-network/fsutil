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
	Object() Object
	Timestamp() time.Time
	TypeMask() EventTypeMask
}
