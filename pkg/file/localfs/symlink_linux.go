// +build linux

package localfs

import (
	"github.com/my-network/fsutil/pkg/file"
	"golang.org/x/sys/unix"
)

const (
	PathMax = 4096 // see 'grep PATH_MAX /usr/include/linux/limits.h'
)

func (symlink *Symlink) Destination() (file.Path, error) {
	destinationPathBuf := make([]byte, PathMax)
	n, err := unix.Readlinkat(int(symlink.Backend.Fd()), "", destinationPathBuf)
	if err != nil {
		return nil, err
	}
	return localToPathBytes(destinationPathBuf[:n]), nil
}
