package localfs

import (
	"io"

	"github.com/my-network/fsutil/pkg/file"
)

var _ file.Socket = &Socket{}

type Socket struct {
	Object
}

func (sock *Socket) Connect() (io.ReadWriteCloser, error) {
	return nil, file.ErrNotImplemented{}
}
