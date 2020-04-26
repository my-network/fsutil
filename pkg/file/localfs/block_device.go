package localfs

import (
	"github.com/my-network/fsutil/pkg/file"
)

var _ file.Device = &BlockDevice{}

type BlockDevice struct {
	Object
}

func (dev *BlockDevice) Read(b []byte) (int, error) {
	return dev.Backend.Read(b)
}

func (dev *BlockDevice) Write(b []byte) (int, error) {
	return dev.Backend.Write(b)
}

func (dev *BlockDevice) ReadAt(b []byte, offset int64) (int, error) {
	return dev.Backend.ReadAt(b, offset)
}

func (dev *BlockDevice) WriteAt(b []byte, offset int64) (int, error) {
	return dev.Backend.WriteAt(b, offset)
}

func (dev *BlockDevice) Seek(offset int64, whence int) (int64, error) {
	return dev.Backend.Seek(offset, whence)
}

func (dev *BlockDevice) IsCharacter() bool {
	return false
}

func (dev *BlockDevice) DevID() interface{} {
	panic("not implemented, yet")
}
