package file

import (
	"fmt"
	"os"
)

type ErrNotImplemented struct{}

func (err ErrNotImplemented) Error() string {
	return "not implemented, yet"
}

type ErrAborted struct{}

func (err ErrAborted) Error() string {
	return "aborted"
}

type ErrWatch struct {
	Path Path
	Err  error
}

func (err ErrWatch) Error() string {
	return fmt.Sprintf("unable to start watching '%s': %v",
		err.Path.LocalPath(), err.Err)
}

type ErrGetChildrenInfo struct {
	Dir Directory
	Err error
}

func (err ErrGetChildrenInfo) Error() string {
	return fmt.Sprintf("unable to get children of path '%s': %v", err.Dir.Path().LocalPath(), err.Err)
}

func (err ErrGetChildrenInfo) Unwrap() error {
	return err.Err
}

type ErrWalkCallback struct {
	Dir   Directory
	Child os.FileInfo
	Err   error
}

func (err ErrWalkCallback) Error() string {
	return fmt.Sprintf("got error from callback on '%s' / '%s': %v",
		err.Dir.Path().LocalPath(), err.Child.Name(), err.Err)
}

func (err ErrWalkCallback) Unwrap() error {
	return err.Err
}

type ErrWalkOpen struct {
	Dir   Directory
	Child os.FileInfo
	Err   error
}

func (err ErrWalkOpen) Error() string {
	if err.Dir == nil {
		return fmt.Sprintf("cannot open the root: %v",
			err.Err)
	}
	return fmt.Sprintf("cannot open '%s' of '%s': %v",
		err.Child.Name(), err.Dir.Path().LocalPath(), err.Err)
}

func (err ErrWalkOpen) Unwrap() error {
	return err.Err
}

type ErrWalkNotDir struct {
	Dir   Directory
	Child Object
}

func (err ErrWalkNotDir) Error() string {
	if err.Dir == err.Child {
		return fmt.Sprintf("root '%s' is not a directory: %T",
			err.Dir.Path().LocalPath(), err.Child)
	}
	return fmt.Sprintf("child '%s' of '%s' is not a directory: %T",
		err.Child.Name(), err.Dir.Path().LocalPath(), err.Child)
}

type ErrWatchMark struct {
	Path Path
	Err  error
}

func (err ErrWatchMark) Error() string {
	return fmt.Sprintf("unable to mark '%s' to be watched: %v",
		err.Path.LocalPath(), err.Err)
}

func (err ErrWatchMark) Unwrap() error {
	return err.Err
}
