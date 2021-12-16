// +build linux freebsd darwin

package localfs

import (
	"syscall"
)

func (obj *Object) ID() interface{} {
	return obj.IDUNIX()
}

func (obj *Object) IDUNIX() ObjectIDUNIX {
	st := obj.LastInfo.Sys().(*syscall.Stat_t)
	return ObjectIDUNIX{
		Dev: st.Dev,
		Ino: st.Ino,
	}
}
