package utils

import (
	"os"
	"path/filepath"
	"syscall"
)

func Openat(dirFD uintptr, path string, flag int, perm os.FileMode) (*os.File, error) {
	fd, err := syscall.Openat(int(dirFD), path, flag, uint32(perm))
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: path, Err: err}
	}
	return os.NewFile(uintptr(fd), filepath.Base(path)), nil
}
