package localfs

import (
	"github.com/my-network/fsutil/pkg/file"
)

var _ file.Object = &Untyped{}

type Untyped struct {
	Object
}
