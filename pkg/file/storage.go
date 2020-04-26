package file

import (
	"context"
	"os"
	"time"
)

type Storage interface {
	Open(ctx context.Context, dirAt Object, path Path, mask OpenFlag, defaultPerm os.FileMode) (Object, error)
	Stat(ctx context.Context, dirAt Object, path Path, noFollow bool) (os.FileInfo, error)

	Symlink(ctx context.Context, dirAt Object, path Path, destination Path) error
	Readlink(ctx context.Context, dirAt Object, path Path) (Path, error)

	Mkdir(ctx context.Context, dirAt Object, path Path, perms os.FileMode, isRecursive bool) error
	Remove(ctx context.Context, dirAt Object, path Path, isRecursive bool) error
	Rename(ctx context.Context, dirAt Object, path, newPath Path) error

	Link(ctx context.Context, dirAt Object, path, destination Path) error

	Chmod(ctx context.Context, dirAt Object, path Path, mode os.FileMode) error
	Chown(ctx context.Context, dirAt Object, path Path, uid, gid int, noFollow bool) error
	Chtimes(ctx context.Context, dirAt Object, path Path, atime time.Time, mtime time.Time) error

	ToLocalPath(Path) string
}

type StorageWatchable interface {
	Storage
	Watcher
}
