package localfs

import (
	"time"

	"github.com/my-network/fsutil/pkg/file"
)

var _ file.Event = &Event{}

type Event struct {
	ObjectValue    file.Object
	TypeMaskValue  file.EventTypeMask
	TimestampValue time.Time
}

func (ev *Event) Object() file.Object {
	return ev.ObjectValue
}

func (ev *Event) TypeMask() file.EventTypeMask {
	return ev.TypeMaskValue
}

func (ev *Event) Timestamp() time.Time {
	return ev.TimestampValue
}
