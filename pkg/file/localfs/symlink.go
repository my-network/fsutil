package localfs

import (
	"context"
	"fmt"
	"os"

	"github.com/my-network/fsutil/pkg/file"
)

var _ file.SymLink = &Symlink{}

type Symlink struct {
	Object
}

func (symlink *Symlink) Open(ctx context.Context, flags file.OpenFlag, defaultPerms os.FileMode) (file.Object, error) {
	destination, err := symlink.Destination()
	if err != nil {
		return nil, fmt.Errorf("unable to get symlink destination: %w", err)
	}
	return symlink.StorageValue.Open(ctx, nil, symlink.Path().Up().Append(destination...), flags, defaultPerms)
}
