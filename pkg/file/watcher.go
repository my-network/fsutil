package file

import (
	"os"
)

type ShouldWatchFunc func(Directory, os.FileInfo) bool

type Watcher interface {
	Watch(Path, ShouldWatchFunc, ShouldWalkFunc) (EventEmitter, error)
}
