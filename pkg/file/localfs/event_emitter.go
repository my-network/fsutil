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
) error {
	if dirAt != nil {
		return file.ErrNotImplemented{}
	}

	var watchErrors []error

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
			pathLocal := dir.Storage().ToLocalPath(dir.Path().Append(objectInfo.Name()))
			err := evEmitter.watcher.Watch(pathLocal)
			fmt.Printf("MARK %v -> %v\n", pathLocal, err)
			if err != nil {
				watchErrors = append(watchErrors, fmt.Errorf("unable to mark '%s' to be watched: %w",
					pathLocal, err))
			}
			return nil
		},
		shouldWalkFunc,
	)

	if err != nil || len(watchErrors) != 0 {
		return &file.ErrCannotWatch{
			Path:        path,
			WalkError:   err,
			WatchErrors: watchErrors,
		}
	}

	return nil
}
