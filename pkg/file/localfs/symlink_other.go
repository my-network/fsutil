// +build !linux

package localfs

import (
	"github.com/my-network/fsutil/pkg/file"
)

func (symlink *Symlink) Destination() (file.Path, error) {
	return nil, file.ErrNotImplemented{}
}
