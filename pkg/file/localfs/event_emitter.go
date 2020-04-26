package localfs

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/howeyc/fsnotify"
	"github.com/my-network/fsutil/pkg/file"
)

var _ file.EventEmitter = &EventEmitter{}

type EventEmitter struct {
	ctx       context.Context
	cancelFn  context.CancelFunc
	storage   *Storage
	watcher   *fsnotify.Watcher
	wg        sync.WaitGroup
	eventChan chan file.Event
}

func newEventEmitter(ctx context.Context, storage *Storage, watcher *fsnotify.Watcher) *EventEmitter {
	evEmitter := &EventEmitter{
		storage:   storage,
		watcher:   watcher,
		eventChan: make(chan file.Event, 1<<16),
	}
	evEmitter.ctx, evEmitter.cancelFn = context.WithCancel(ctx)
	evEmitter.initPipeline()
	return evEmitter
}

func (evEmitter *EventEmitter) initPipeline() {
	evEmitter.wg.Add(1)
	go func() {
		defer func() {
			close(evEmitter.eventChan)
			evEmitter.wg.Done()
		}()
		evEmitter.pipelineLoop()
	}()
}

func (evEmitter *EventEmitter) pipelineLoop() {
	for {
		select {
		case ev := <-evEmitter.watcher.Event:
			now := time.Now()

			evEmitter.eventChan <- &Event{
				PathValue:      localToPath(ev.Name).RelativeTo(evEmitter.storage.workDir),
				TypeMaskValue:  0,
				TimestampValue: now,
			}
		case <-evEmitter.ctx.Done():
		}
	}
}

func (evEmitter *EventEmitter) C() <-chan file.Event {
	return evEmitter.eventChan
}

func (evEmitter *EventEmitter) Close() error {
	evEmitter.cancelFn()
	evEmitter.wg.Wait()
	return nil
}

func (evEmitter *EventEmitter) Watch(
	dirAt file.Directory,
	path file.Path,
	shouldWatchFunc file.ShouldWatchFunc,
	shouldWalkFunc file.ShouldWalkFunc,
	errorHandler file.ErrorHandlerFunc,
) error {
	if errorHandler == nil {
		errorHandler = dummyErrorHandler
	}

	if dirAt != nil {
		return file.ErrNotImplemented{}
	}

	err := file.Walk(
		evEmitter.ctx,
		evEmitter.storage,
		nil,
		path,
		func(dir file.Directory, objectInfo os.FileInfo) error {
			fmt.Printf("MARK? %v %v %v\n", dir.Path().LocalPath(), objectInfo.Name(), objectInfo.IsDir())
			if !objectInfo.IsDir() {
				return nil
			}
			if shouldWatchFunc != nil && !shouldWatchFunc(dir, objectInfo) {
				return nil
			}
			pathFull := dir.Path().Append(objectInfo.Name())
			pathFullLocal := dir.Storage().ToLocalPath(pathFull)
			err := evEmitter.watcher.Watch(pathFullLocal)
			fmt.Printf("MARK %v -> %v\n", pathFullLocal, err)
			if err != nil {
				if err := errorHandler(file.ErrWatchMark{Path: pathFull, Err: err}); err != nil {
					return err
				}
			}
			return nil
		},
		shouldWalkFunc,
		errorHandler,
	)

	if err != nil {
		return &file.ErrWatch{
			Path: path,
			Err:  err,
		}
	}

	return nil
}
