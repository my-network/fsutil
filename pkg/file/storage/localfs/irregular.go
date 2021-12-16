package localfs

import (
	"github.com/my-network/fsutil/pkg/file"
)

var _ file.Object = &Irregular{}

type Irregular struct {
	Object
}
