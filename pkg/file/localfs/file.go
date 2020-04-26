package localfs

import (
	"github.com/my-network/fsutil/pkg/file"
)

var _ file.File = &File{}

type File struct {
	Object
}

func (f *File) Read(b []byte) (int, error) {
	return f.Backend.Read(b)
}

func (f *File) Write(b []byte) (int, error) {
	return f.Backend.Write(b)
}

func (f *File) ReadAt(b []byte, offset int64) (int, error) {
	return f.Backend.ReadAt(b, offset)
}

func (f *File) WriteAt(b []byte, offset int64) (int, error) {
	return f.Backend.WriteAt(b, offset)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.Backend.Seek(offset, whence)
}
