package localfs

import (
	"os"

	"github.com/my-network/fsutil/pkg/file"
)

type Backend = *os.File

type ObjectIDUNIX struct {
	Dev uint64
	Ino uint64
}

type Object struct {
	Backend
	StorageValue *Storage
	LastInfo     os.FileInfo
	LastPath     file.Path
}

func (obj *Object) Name() string {
	return obj.Backend.Name()
}

func (obj *Object) Stat() (os.FileInfo, error) {
	newInfo, err := obj.Backend.Stat()
	if err != nil {
		return nil, err
	}
	obj.LastInfo = newInfo
	return newInfo, nil
}

func (obj *Object) LastStat() os.FileInfo {
	return obj.LastInfo
}

func (obj *Object) Path() file.Path {
	// TODO: validate if the path is still valid
	return obj.LastPath
}

func (obj *Object) Close() error {
	return obj.Backend.Close()
}

func (obj *Object) Storage() file.Storage {
	return obj.StorageValue
}

func (obj *Object) FD() uintptr {
	return obj.Backend.Fd()
}
