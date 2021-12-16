package event

import (
	"os"

	"github.com/my-network/fsutil/pkg/file"
)

type ShouldWatchFunc func(file.Directory, os.FileInfo) bool

type Watcher interface {
	Watch(file.Directory, file.Path, ShouldWatchFunc, file.ShouldWalkFunc, file.ErrorHandlerFunc) (Emitter, error)
}
