package event

import (
	"github.com/my-network/fsutil/pkg/file"
)

type Emitter interface {
	Close() error
	C() <-chan Event
	Watch(file.Directory, file.Path, ShouldWatchFunc, file.ShouldWalkFunc, file.ErrorHandlerFunc) error
}
