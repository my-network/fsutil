package cached

import (
	"context"
	"os"
	"path/filepath"

	"github.com/my-network/fsutil/pkg/file"
	"github.com/xaionaro-go/synctools"
)

type Map struct {
	synctools.RWMutex
	Map map[string]*Object
}

type Storage struct {
	file.Storage
	Map
	Config
}

func NewStorage(backendStorage file.Storage, opts ...Option) *Storage {
	storage := &Storage{
		Storage: backendStorage,
		Map: Map{
			Map: map[string]*Object{},
		},
	}
	for _, opt := range opts {
		opt.apply(&storage.Config)
	}
	return storage
}

func (stor *Storage) OpenInBackground(
	dirAt file.Object,
	path file.Path,
	mask file.OpenFlag,
	defaultPerm os.FileMode,
) {
	stor.openInBackground(context.Background(), dirAt, path, mask, defaultPerm)
}

func (stor *Storage) openInBackground(
	ctx context.Context,
	dirAt file.Object,
	path file.Path,
	mask file.OpenFlag,
	defaultPerm os.FileMode,
) *Object {
	if dirAt != nil {
		panic(file.ErrNotImplemented{})
	}

	var obj *Object
	stor.Map.LockDo(func() {
		mapKey := filepath.Join(path...)
		obj = stor.Map.Map[mapKey]
		if obj != nil {
			return
		}

		obj = &Object{}
		stor.Map.Map[mapKey] = obj

		obj.Lock()
		go func() {
			defer obj.Unlock()

			obj.Object, obj.OpenError = stor.Storage.Open(
				ctx,
				dirAt,
				path,
				mask,
				defaultPerm,
			)
		}()
	})
	return obj
}

func (stor *Storage) Open(
	ctx context.Context,
	dirAt file.Object,
	path file.Path,
	mask file.OpenFlag,
	defaultPerm os.FileMode,
) (file.Object, error) {
	obj := stor.openInBackground(ctx, dirAt, path, mask, defaultPerm)
	obj.RLock()
	defer obj.RUnlock()
	return obj, obj.OpenError
}

func (stor *Storage) OpenCached(path []string) (file.Object, error) {
	var obj *Object
	stor.Map.RLockDo(func() {
		mapKey := filepath.Join(path...)
		obj = stor.Map.Map[mapKey]
	})
	if obj == nil {
		return nil, ErrNotOpened{Path: path}
	}
	return obj, obj.OpenError
}
