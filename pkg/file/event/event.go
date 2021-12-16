package event

import (
	"time"

	pkgbytes "github.com/my-network/fsutil/pkg/bytes"
	"github.com/my-network/fsutil/pkg/file"
)

type Event struct {
	ObjID     interface{}
	Path      file.Path
	TypeMask  TypeMask
	Timestamp time.Time

	Range   *pkgbytes.Range
	MovedTo file.Path
}
