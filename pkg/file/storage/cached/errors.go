package cached

import (
	"fmt"
)

type ErrNotOpened struct {
	Path []string
}

func (err ErrNotOpened) Error() string {
	return fmt.Sprintf("object '%v' is not opened", err.Path)
}
