package localfs

import (
	"context"
	"os"

	"github.com/my-network/fsutil/pkg/file"
)

var _ file.PathDescriptor = &PathDescriptor{}

type PathDescriptor struct {
	Object
}

func (pathDesc *PathDescriptor) Open(
	ctx context.Context,
	mask file.OpenFlag,
	defaultPerm os.FileMode,
) (file.Object, error) {
	return nil, file.ErrNotImplemented{}
}
