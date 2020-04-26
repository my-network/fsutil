package file

import (
	"os"
)

type ShouldWatchFunc func(Directory, os.FileInfo) bool

type Watcher interface {
	Watch(Directory, Path, ShouldWatchFunc, ShouldWalkFunc, ErrorHandlerFunc) (EventEmitter, error)
}
