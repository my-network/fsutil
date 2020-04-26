package localfs

import (
	"github.com/my-network/fsutil/pkg/file"
)

var _ file.Pipe = &NamedPipe{}

type NamedPipe struct {
	Object
}

func (pipe *NamedPipe) Read(b []byte) (int, error) {
	return pipe.Backend.Read(b)
}

func (pipe *NamedPipe) Write(b []byte) (int, error) {
	return pipe.Backend.Write(b)
}
