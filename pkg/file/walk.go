package file

import (
	"context"
	"fmt"
	"os"
)

type CallbackFunc func(dir Directory, file os.FileInfo) error
type ShouldWalkFunc func(dir Directory, file os.FileInfo) bool

func Walk(
	ctx context.Context,
	storage Storage,
	root Path,
	callback CallbackFunc,
	shouldWalkFn ShouldWalkFunc,
) error {
	dirObj, err := storage.Open(ctx, nil, root, FlagWalkDefaults, 0000)
	if err != nil {
		return fmt.Errorf("unable to open root '%s': %w",
			root.LocalPath(), err)
	}

	dir, ok := dirObj.(Directory)
	if !ok {
		return fmt.Errorf("root '%s' is not a directory: %T",
			root.LocalPath(), dirObj)
	}

	return walkDir(ctx, dir, callback, shouldWalkFn)
}

func walkDir(
	ctx context.Context,
	dir Directory,
	callback CallbackFunc,
	shouldWalkFn ShouldWalkFunc,
) error {
	children, err := dir.Readdir(-1)
	if err != nil {
		return fmt.Errorf("unable to get children of path '%s': %w", dir.Path().LocalPath(), err)
	}
	for _, childInfo := range children {
		err = callback(dir, childInfo)
		if err != nil {
			return fmt.Errorf("got error from callback on '%s': %w", dir.Path().LocalPath(), err)
		}

		if !childInfo.IsDir() {
			continue
		}
		if shouldWalkFn != nil && !shouldWalkFn(dir, childInfo) {
			continue
		}

		childObj, err := dir.Open(ctx, Path{childInfo.Name()}, FlagWalkDefaults, 0000)
		if err != nil {
			return fmt.Errorf("unable to open child '%s' of '%s': %w",
				childInfo.Name(), dir.Path().LocalPath(), err)
		}

		child, ok := childObj.(Directory)
		if !ok {
			return fmt.Errorf("child '%s' of '%s' is not a directory: %T",
				childInfo.Name(), dir.Path().LocalPath(), childObj)
		}

		err = walkDir(ctx, child, callback, shouldWalkFn)
		if err != nil {
			return err
		}
	}

	return nil
}
