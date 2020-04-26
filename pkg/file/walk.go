package file

import (
	"context"
	"os"
)

type CallbackFunc func(dir Directory, file os.FileInfo) error
type ShouldWalkFunc func(dir Directory, file os.FileInfo) bool
type ErrorHandlerFunc func(error) error

type curDirInfo struct {
	os.FileInfo
}

func (info curDirInfo) Name() string {
	return "."
}

func dummyErrorHandler(err error) error {
	return err
}

func Walk(
	ctx context.Context,
	storage Storage,
	dirAt Directory,
	path Path,
	callback CallbackFunc,
	shouldWalkFn ShouldWalkFunc,
	errorHandlerFn ErrorHandlerFunc,
) error {
	if errorHandlerFn == nil {
		errorHandlerFn = dummyErrorHandler
	}

	dirObj, err := storage.Open(ctx, dirAt, path, FlagWalkDefaults, 0000)
	if err != nil {
		if err := errorHandlerFn(ErrWalkOpen{Dir: nil, Child: nil, Err: err}); err != nil {
			return err
		}
	}

	dir, ok := dirObj.(Directory)
	if !ok {
		if err := errorHandlerFn(ErrWalkNotDir{Dir: dir, Child: dir}); err != nil {
			return err
		}
	}

	dirInfo := curDirInfo{FileInfo: dir.LastStat()}
	err = callback(dir, dirInfo)
	if err != nil {
		if err := errorHandlerFn(ErrWalkCallback{Dir: dir, Child: dirInfo, Err: err}); err != nil {
			return err
		}
	}

	return walkDir(ctx, dir, callback, shouldWalkFn, errorHandlerFn)
}

func walkDir(
	ctx context.Context,
	dir Directory,
	callback CallbackFunc,
	shouldWalkFn ShouldWalkFunc,
	errorHandlerFn ErrorHandlerFunc,
) error {
	children, err := dir.Readdir(-1)
	if err != nil {
		if err = errorHandlerFn(ErrGetChildrenInfo{Dir: dir, Err: err}); err != nil {
			return err
		}
	}
	for _, childInfo := range children {
		err = callback(dir, childInfo)
		if err != nil {
			if err := errorHandlerFn(ErrWalkCallback{Dir: dir, Child: childInfo, Err: err}); err != nil {
				return err
			}
		}

		if !childInfo.IsDir() {
			continue
		}
		if shouldWalkFn != nil && !shouldWalkFn(dir, childInfo) {
			continue
		}

		childObj, err := dir.Open(ctx, Path{childInfo.Name()}, FlagWalkDefaults, 0000)
		if err != nil {
			if err := errorHandlerFn(ErrWalkOpen{Dir: dir, Child: childInfo, Err: err}); err != nil {
				return err
			}
		}

		child, ok := childObj.(Directory)
		if !ok {
			if err := errorHandlerFn(ErrWalkNotDir{Dir: dir, Child: childObj}); err != nil {
				return err
			}
		}

		err = walkDir(ctx, child, callback, shouldWalkFn, errorHandlerFn)
		if err != nil {
			return err
		}
	}

	return nil
}
