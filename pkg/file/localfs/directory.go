package localfs

import (
	"context"
	"os"

	"github.com/my-network/fsutil/pkg/file"
)

var _ file.Directory = &Directory{}

type Directory struct {
	Object
}

func (dir *Directory) Readdir(n int) ([]os.FileInfo, error) {
	return dir.Backend.Readdir(n)
}

func (dir *Directory) Open(
	ctx context.Context,
	path file.Path,
	flags file.OpenFlag,
	defaultPerm os.FileMode,
) (file.Object, error) {
	return dir.StorageValue.Open(ctx, dir, path, flags, defaultPerm)
}
