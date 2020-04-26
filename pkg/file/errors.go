package file

import (
	"fmt"
)

type ErrNotImplemented struct{}

func (err ErrNotImplemented) Error() string {
	return "not implemented, yet"
}

type ErrAborted struct{}

func (err ErrAborted) Error() string {
	return "aborted"
}

type ErrCannotWatch struct {
	Path        Path
	WalkError   error
	WatchErrors []error
}

func (err ErrCannotWatch) Error() string {
	return fmt.Sprintf("unable to start watching '%s': walkErr:%v; watchErrs: %v",
		err.Path.LocalPath(), err.WalkError, err.WatchErrors)
}
