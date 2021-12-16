package localfs

import (
	"github.com/my-network/fsutil/pkg/file"
)

var _ file.Device = &CharDevice{}

type CharDevice struct {
	Object
}

func (dev *CharDevice) Read(b []byte) (int, error) {
	return dev.Backend.Read(b)
}

func (dev *CharDevice) Write(b []byte) (int, error) {
	return dev.Backend.Write(b)
}

func (dev *CharDevice) ReadAt(b []byte, offset int64) (int, error) {
	return dev.Backend.ReadAt(b, offset)
}

func (dev *CharDevice) WriteAt(b []byte, offset int64) (int, error) {
	return dev.Backend.WriteAt(b, offset)
}

func (dev *CharDevice) Seek(offset int64, whence int) (int64, error) {
	return dev.Backend.Seek(offset, whence)
}

func (dev *CharDevice) IsCharacter() bool {
	return true
}

func (dev *CharDevice) DevID() interface{} {
	panic("not implemented, yet")
}
