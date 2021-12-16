package cached

import (
	"sync"

	"github.com/my-network/fsutil/pkg/file"
)

type Object struct {
	sync.RWMutex
	file.Object
	OpenError error
}
