package file

import (
	"context"
	"io"
	"os"
	"time"
)

type Object interface {
	Path() Path
	Name() string

	Stat() (os.FileInfo, error)
	LastStat() os.FileInfo

	Close() error
	Chown(uid, gid int) error

	Storage() Storage
	FD() uintptr
}

type PathDescriptor interface {
	Object
	Chmod(mode os.FileMode) error

	Open(context.Context, OpenFlag, os.FileMode) (Object, error)
}

type File interface {
	Object
	IO
	Chmod(mode os.FileMode) error

	Seek(offset int64, whence int) (int64, error)

	WriteAt(b []byte, off int64) (n int, err error)
	ReadAt(b []byte, off int64) (n int, err error)
}

type IO interface {
	Write(b []byte) (n int, err error)
	Read(b []byte) (n int, err error)

	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error

	Sync() error
}

type Device interface {
	File

	IsCharacter() bool
	DevID() interface{}
}

type Pipe interface {
	Object
	Chmod(mode os.FileMode) error
	IO
}

type Socket interface {
	Object
	Chmod(mode os.FileMode) error

	Connect() (io.ReadWriteCloser, error)
}

type SymLink interface {
	PathDescriptor

	Destination() (Path, error)
}

type Directory interface {
	Object
	Chmod(mode os.FileMode) error

	Readdir(int) ([]os.FileInfo, error)
	Open(context.Context, Path, OpenFlag, os.FileMode) (Object, error)
}
