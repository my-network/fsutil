package localfs

import (
	"time"

	"github.com/my-network/fsutil/pkg/file"
)

var _ file.Event = &Event{}

type Event struct {
	PathValue      file.Path
	TypeMaskValue  file.EventTypeMask
	TimestampValue time.Time
}

func (ev *Event) Path() file.Path {
	return ev.PathValue
}

func (ev *Event) TypeMask() file.EventTypeMask {
	return ev.TypeMaskValue
}

func (ev *Event) Timestamp() time.Time {
	return ev.TimestampValue
}
