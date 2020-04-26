package localfs

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/howeyc/fsnotify"
	"github.com/my-network/fsutil/pkg/file"
	"github.com/my-network/fsutil/pkg/file/utils"
)

var _ file.StorageWatchable = &Storage{}

type Storage struct {
	ctx      context.Context
	cancelFn context.CancelFunc
	wg       sync.WaitGroup
	workDir  file.Path
}

func NewStorage(workDir string) *Storage {
	stor := &Storage{
		workDir: file.ParseLocalPath(workDir),
	}
	stor.ctx, stor.cancelFn = context.WithCancel(context.Background())
	return stor
}

func (stor *Storage) WorkDir() file.Path {
	result := make(file.Path, len(stor.workDir))
	copy(result, stor.workDir)
	return result
}

func (stor *Storage) Watch(
	dirAt file.Directory,
	path file.Path,
	shouldMarkFunc file.ShouldWatchFunc,
	shouldWalkFunc file.ShouldWalkFunc,
) (file.EventEmitter, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize a watcher backend: %w", err)
	}

	evEmitter := newEventEmitter(stor.ctx, stor, watcher)

	err = evEmitter.Watch(dirAt, path, shouldMarkFunc, shouldWalkFunc)
	if err != nil {
		return nil, err
	}

	stor.wg.Add(1)
	go func() {
		defer stor.wg.Done()
		select {
		case <-stor.ctx.Done():
		}
		_ = evEmitter.Close()
	}()

	return evEmitter, nil
}

func (stor *Storage) Close() error {
	stor.cancelFn()
	stor.wg.Wait()
	return nil
}

func (stor *Storage) ToAbsPath(pathRel file.Path) file.Path {
	return stor.workDir.Append(pathRel...)
}

func (stor *Storage) ToLocalPath(path file.Path) string {
	return stor.ToAbsPath(path).LocalPath()
}

func (stor *Storage) ToLocalPathAt(dir file.Object, path file.Path) string {
	if dir != nil {
		return path.LocalPath()
	}
	return stor.ToLocalPath(path)
}

func (stor *Storage) Open(
	ctx context.Context,
	dirAt file.Object,
	path file.Path,
	flags file.OpenFlag,
	defaultPerm os.FileMode,
) (file.Object, error) {
	select {
	case <-ctx.Done():
		return nil, file.ErrAborted{}
	default:
	}

	var f *os.File
	var err error
	if dirAt != nil {
		f, err = utils.Openat(dirAt.FD(), path.LocalPath(), flags.OSFlags(), defaultPerm)
	} else {
		f, err = os.OpenFile(stor.ToLocalPath(path), flags.OSFlags(), defaultPerm)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to open '%s': %w",
			path.LocalPath(), err)
	}

	select {
	case <-ctx.Done():
		return nil, file.ErrAborted{}
	default:
	}

	fileInfo, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to 'stat' on '%s': %w",
			path.LocalPath(), err)
	}

	obj := Object{
		StorageValue: stor,
		Backend:      f,
		LastInfo:     fileInfo,
		LastPath:     path,
	}

	switch fileInfo.Mode() & os.ModeType {
	case os.ModeDir:
		return &Directory{Object: obj}, nil
	case os.ModeSymlink:
		return &Symlink{Object: obj}, nil
	case os.ModeSocket:
		return &Socket{Object: obj}, nil
	case os.ModeIrregular:
		return &Irregular{Object: obj}, nil
	}

	if flags.HasPath() {
		switch fileInfo.Mode() & os.ModeType {
		case 0, os.ModeNamedPipe, os.ModeDevice, os.ModeDevice | os.ModeCharDevice:
			return &PathDescriptor{Object: obj}, nil
		}
	} else {
		switch fileInfo.Mode() & os.ModeType {
		case 0:
			return &File{Object: obj}, nil
		case os.ModeDevice:
			return &BlockDevice{Object: obj}, nil
		case os.ModeDevice | os.ModeCharDevice:
			return &CharDevice{Object: obj}, nil
		case os.ModeNamedPipe:
			return &NamedPipe{Object: obj}, nil
		}
	}
	return &Untyped{Object: obj}, nil

}

func (stor *Storage) Stat(
	ctx context.Context,
	dirAt file.Object,
	path file.Path,
	noFollow bool,
) (os.FileInfo, error) {
	if dirAt != nil {
		return nil, file.ErrNotImplemented{}
	}

	var statFunc func(name string) (os.FileInfo, error)
	if noFollow {
		statFunc = os.Lstat
	} else {
		statFunc = os.Stat
	}
	return statFunc(stor.ToLocalPath(path))
}

func (stor *Storage) Symlink(
	ctx context.Context,
	dirAt file.Object,
	path file.Path,
	destination file.Path,
) error {
	if dirAt != nil {
		return file.ErrNotImplemented{}
	}

	return os.Symlink(destination.LocalPath(), stor.ToLocalPath(path))
}

func (stor *Storage) Readlink(
	ctx context.Context,
	dirAt file.Object,
	path file.Path,
) (file.Path, error) {
	if dirAt != nil {
		return nil, file.ErrNotImplemented{}
	}

	destination, err := os.Readlink(stor.ToLocalPath(path))
	if err != nil {
		return nil, err
	}
	return localToPath(destination), nil
}

func (stor *Storage) Mkdir(
	ctx context.Context,
	dirAt file.Object,
	path file.Path,
	perms os.FileMode,
	recursive bool,
) error {
	if dirAt != nil {
		return file.ErrNotImplemented{}
	}

	var mkdirFunc func(name string, mode os.FileMode) error
	if recursive {
		mkdirFunc = os.MkdirAll
	} else {
		mkdirFunc = os.Mkdir
	}

	return mkdirFunc(stor.ToLocalPath(path), perms)
}

func (stor *Storage) Remove(
	ctx context.Context,
	dirAt file.Object,
	path file.Path,
	isRecursive bool,
) error {
	if dirAt != nil {
		return file.ErrNotImplemented{}
	}

	var removeFunc func(name string) error
	if isRecursive {
		removeFunc = os.RemoveAll
	} else {
		removeFunc = os.Remove
	}

	return removeFunc(stor.ToLocalPath(path))
}

func (stor *Storage) Rename(
	ctx context.Context,
	dirAt file.Object,
	path, newPath file.Path,
) error {
	if dirAt != nil {
		return file.ErrNotImplemented{}
	}

	return os.Rename(stor.ToLocalPath(path), stor.ToLocalPath(newPath))
}

func (stor *Storage) Link(
	ctx context.Context,
	dirAt file.Object,
	path, destination file.Path,
) error {
	if dirAt != nil {
		return file.ErrNotImplemented{}
	}

	return os.Link(stor.ToLocalPath(path), stor.ToLocalPath(destination))
}

func (stor *Storage) Chmod(
	ctx context.Context,
	dirAt file.Object,
	path file.Path,
	mode os.FileMode,
) error {
	if dirAt != nil {
		return file.ErrNotImplemented{}
	}

	return os.Chmod(stor.ToLocalPath(path), mode)
}

func (stor *Storage) Chown(
	ctx context.Context,
	dirAt file.Object,
	path file.Path,
	uid, gid int,
	noFollow bool,
) error {
	if dirAt != nil {
		return file.ErrNotImplemented{}
	}

	var chownFunc func(name string, uid, gid int) error
	if noFollow {
		chownFunc = os.Lchown
	} else {
		chownFunc = os.Chown
	}

	return chownFunc(stor.ToLocalPath(path), uid, gid)
}

func (stor *Storage) Chtimes(
	ctx context.Context,
	dirAt file.Object,
	path file.Path,
	atime, mtime time.Time,
) error {
	if dirAt != nil {
		return file.ErrNotImplemented{}
	}

	return os.Chtimes(stor.ToLocalPath(path), atime, mtime)
}
